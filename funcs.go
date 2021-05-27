package main

import (
	"bufio"
	"flag"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
	"fmt"
)


func usage() {
	fmt.Fprintf(os.Stderr, `Version: k8slog v1.0.0
Tail multiple pods and containers logs from Kubernetes

Usage: ./k8slog  [-n namespace] [-f tail lines] [-s since ] [-l k8s label Selector] [-c container name] [-r running inside or outside k8s,default is "outside" k8s] [-kubeconfig  absolute path to the kubeconfig file]

Examples:
./k8slog -n mynamespace -l label=singlecontainer -f 100  
./k8slog -n mynamespace -l label=muticontainer -f 100 -c containername
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -r inside
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -kubeconfig admin.kubeconfig
Options:
`)
	flag.PrintDefaults()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func onaddpod(obj interface{}) {
}

func onupdatepod(old,new interface{})  {

	newpod := new.(*corev1.Pod)

	if newpod.Status.Phase=="Running" {
		mutex.Lock()
		defer mutex.Unlock()
		var containerstatus corev1.ContainerStatus

		if subcontainer!=""{
			containerstatus=Getcontainer(newpod,subcontainer)
		}else {
			containerstatus=newpod.Status.ContainerStatuses[0]
		}

		if containerstatus.Ready==true && ! Existinpodlist(targetlistpods,newpod.Name){
			log.Println(newpod.Name,containerstatus.Name," pod status is ready ,可以打印日志了",)
			targetlistpods=append(targetlistpods,*newpod)
			newpodchan<-*newpod

		}else if containerstatus.Ready==false && Existinpodlist(targetlistpods,newpod.Name) {
			targetlistpods=Removepodlist(targetlistpods,newpod.Name)
		}
	}else {
		if Existinpodlist(targetlistpods,newpod.Name){
			targetlistpods=Removepodlist(targetlistpods,newpod.Name)
		}
	}

}
func ondeletepod(obj interface{}){

	pod := obj.(*corev1.Pod)
	fmt.Println("delete a pod:", pod.Name)
	mutex.Lock()
	defer mutex.Unlock()
	targetlistpods=Removepodlist(targetlistpods,pod.Name)

}

func getpodlist(clientset *kubernetes.Clientset,namespace string) ([]corev1.Pod,error){
	var pods *corev1.PodList

	if podlabel==""{
		pods,err=clientset.CoreV1().Pods(namespace).List(context.Background(),metav1.ListOptions{})

	}else {
		pods,err=clientset.CoreV1().Pods(namespace).List(context.Background(),metav1.ListOptions{LabelSelector: podlabel})

	}
	if err!=nil{
		log.Println("err list pods ",err.Error())
		if err.Error()=="Unauthorized"{
			log.Fatal("Unauthorized error, Please check if you have login to kubernetes cluster!")
		}
		return nil,err
	}

	healthpodlist:=make([]corev1.Pod,0)
	for _,pod:=range pods.Items{
		if pod.Status.Phase=="Running"{
			healthpodlist=append(healthpodlist,pod)
		}
	}

	return healthpodlist,nil
}

func gettaillines(clientset *kubernetes.Clientset,namespace string,pod corev1.Pod)*rest.Request{
	var podlogoptions *corev1.PodLogOptions
	//sincesecond选项优先于tailline选项
	sincesecondint:=sincesecond.Milliseconds()/1000
	if sincesecondint>0{
		tailline=-1
	}

	//对于一个pod中有多个container的情况，强制要带参数-c指定container
	if subcontainer!="" && len(pod.Status.ContainerStatuses)>1 {

		if tailline == -1 {
			if sincesecond>0 {
				podlogoptions = &corev1.PodLogOptions{Follow: true, SinceSeconds: &sincesecondint, Container: subcontainer}
			}
			if sincesecond==0{
				podlogoptions = &corev1.PodLogOptions{Follow: true,Container: subcontainer}

			}
		}else if tailline >-1{
			podlogoptions=&corev1.PodLogOptions{Follow: true,Container:subcontainer,TailLines:&tailline}
		}else {
			log.Println("tailline 非法数字，请重新输入")
		}

	}


	//如果没-c参数指定container，那么默认只打印pod的日志，对于一个pod多个container的情况会报错，要指定container
	if subcontainer=="" {

		if tailline == -1 {
			if sincesecond>0 {
				podlogoptions = &corev1.PodLogOptions{Follow: true, SinceSeconds: &sincesecondint}
			}
			if sincesecond==0{
				podlogoptions = &corev1.PodLogOptions{Follow: true}
			}
		}else if tailline >-1 {
			podlogoptions=&corev1.PodLogOptions{Follow: true,TailLines:&tailline}
		}else {
			log.Println("tailline 非法数字，请重新输入")
		}
	}
	return clientset.CoreV1().Pods(namespace).GetLogs(pod.Name,podlogoptions)
}

func producepodlog(clientset *kubernetes.Clientset,namespace string,pod corev1.Pod)  {
	defer dorecovery()

	//如果参数的-c 不为空，但是pod里面只有一个container，那么直接停止这个pod的日志输出
	//如果-c 为空，对于一个pod里包含多个container的情况，会跳过这个pod的日志输出
	//要打印指定container的日志，必须用-c指定
	if subcontainer!="" && len(pod.Status.ContainerStatuses)<=1{
		return
	}
	if subcontainer=="" && len(pod.Status.ContainerStatuses)>1{
		log.Println("there are",len(pod.Status.ContainerStatuses),"containers inside pod,a container name must be specified for pod ",pod.Name, "choose one of: ",Listcontainer(pod))
		return
	}

	rlog := gettaillines(clientset, namespace, pod)

	logcontent,err:=rlog.Stream(context.Background())
	if err!=nil{
		log.Println("出错了",pod.Name,err)
	}
	defer logcontent.Close()
	r:=bufio.NewReader(logcontent)
	for {
		buf,_,err:=r.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println(pod.Name," 已退出，不会打印日志了 ",err)
				break;
			}
		}
		fmt.Println(pod.Name, string(buf))
	}
}

func dorecovery(){
	if err:=recover();err!=nil{
		log.Println("error happen ",err)
	}
}

func Existinpodlist(podlist []corev1.Pod,podname string) bool {
	for _,pod :=range podlist{
		if pod.Name==podname{
			return true
		}
	}
	return false
}

func Removepodlist(podlist []corev1.Pod,podname string) []corev1.Pod {
	indexpod:=0
	for index,pod :=range podlist{
		if pod.Name==podname{
			indexpod=index
			break
		}
	}
	if indexpod==0{
		return podlist
	}
	return append(podlist[:indexpod],podlist[indexpod+1:]...)
}

func Getcontainer(pod *corev1.Pod,containername string) corev1.ContainerStatus{
	for _,containerstatus:=range pod.Status.ContainerStatuses{
		if containerstatus.Name==containername{
			return containerstatus
		}
	}
	return corev1.ContainerStatus{}
}

func Listcontainer(pod corev1.Pod)(containerlist []string) {

	for _,container :=range pod.Status.ContainerStatuses{
		containerlist=append(containerlist,container.Name)
	}

	return
}

func Checkfileexist(path string)(error){
	_,err:=os.Stat(path)
	return err
}
