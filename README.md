# k8slog
This is a tool to check multiple pods and containers logs from Kubernetes

## Usage: 

./k8slog  [-n namespace] [-f tail lines] [-s since ] [-l k8s label Selector] [-c container name] [-r running inside or outside k8s,default is "outside" k8s] [-kubeconfig  absolute path to the kubeconfig file]

Examples:
./k8slog -n mynamespace -l label=singlecontainer -f 100      
./k8slog -n mynamespace -l label=muticontainer -f 100 -c containername    
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername    
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -r inside    
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -kubeconfig admin.kubeconfig    
## Options:    
  -c string    
        -c container name , Container name when multiple containers in pod    
  -f int    
        -f 10,The number of lines from the end of the logs to show. Defaults to -1, showing all logs. (default -1)    
  -kubeconfig string
        (optional) absolute path to the kubeconfig file,default is ${HOME}/.kube/config    
  -l string    
        -l app=nginx,Selector (label query) to filter on. If present, default to ".*" for the pod-query.    
  -n string
        -n namespace (default "default")    
  -r string
        -r inside means run inside k8s,by default it is running outside k8s (default "outside")    
  -s duration    
        -s 10s : Return logs newer than a relative duration like 5s, 2m, or 2.5h.    
