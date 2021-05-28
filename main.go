package main

import (
	"flag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"path/filepath"
	"sync"
	"time"
)

var (
	subcontainer string
	sincesecond time.Duration
	tailline int64
	podlabel string
	namespace string
	runinsidek8s string
	newpodchan chan corev1.Pod
	targetlistpods []corev1.Pod
	mutex sync.Mutex
	kubeconfig string
	config *rest.Config
	err error

)


func main() {
	flag.Usage=usage

	flag.Int64Var(&tailline,"f",-1,"-f 10,The number of lines from the end of the logs to show. Defaults to -1, showing all logs.")
	flag.DurationVar(&sincesecond,"s",0,"-s 10s : Return logs newer than a relative duration like 5s, 2m, or 2.5h.")
	flag.StringVar(&namespace,"n","default","-n namespace")
	flag.StringVar(&runinsidek8s,"r","outside","-r inside means run inside k8s,by default it is running outside k8s")
	flag.StringVar(&podlabel,"l","","-l app=nginx,Selector (label query) to filter on. If present, default to \".*\" for the pod-query.")
	flag.StringVar(&subcontainer,"c","","-c container name , Container name when multiple containers in pod")
	flag.StringVar(&kubeconfig,"kubeconfig", "", "(optional) absolute path to the kubeconfig file,default is ${HOME}/.kube/config")
	flag.Parse()


	//如果参数少于2个，则打印usage并退出
	if len(flag.Args())>0||flag.CommandLine.NFlag()==0{
		flag.Usage()
		return
	}

	if home := homeDir(); home != ""&&kubeconfig=="" {
		kubeconfig=filepath.Join(home, ".kube", "config")
	}

	switch runinsidek8s {
	case "outside":
		err := Checkfileexist(kubeconfig)
		if err != nil {
			log.Fatal("Get kubeconfig error ", err)
		}
	case "inside":
		kubeconfig=""
	default:
		log.Fatal("runinsidek8s can only be one of [\"inside\",\"outside\"]")

	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//注意，全局变量如果在外侧已经定义，那么不能在函数里再用:=方法来定义，否则这个全局变量会变成局部变量
	targetlistpods,err=getpodlist(clientset,namespace)

	if err!=nil{
		log.Println( err)
	}

	for _,pod:=range targetlistpods{

		go producepodlog(clientset,namespace,pod)
	}


	//死循环来一直监控新的pod
	newpodchan=make(chan corev1.Pod)

	go func() {

		for{
			select {
			case newpod:=<-newpodchan:
				go producepodlog(clientset,namespace,newpod)

			}
		}

	}()



	stopper := make(chan struct{})
	defer close(stopper)

	// 初始化 informer
	factory := informers.NewSharedInformerFactory(clientset, 0)
	podInformer := factory.Core().V1().Pods().Informer()
	defer runtime.HandleCrash()

	// 使用自定义 handler
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onaddpod,
		UpdateFunc: onupdatepod,
		DeleteFunc: ondeletepod,
	})

	// 启动 informer，list & watch
	go factory.Start(stopper)

	<-stopper
}

