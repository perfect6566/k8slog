# k8slog
This is a tool to check multiple pods and containers logs from Kubernetes

## Usage: 

./k8slog  [-n namespace] [-f tail lines] [-s since ] [-l k8s label Selector] [-c container name] [-r running inside or outside k8s,default is "o
utside" k8s] [-kubeconfig  absolute path to the kubeconfig file]

Examples:

//Tail log from pods with label singlecontainer in ns mynamespace and follow the last 100 line logs 
./k8slog -n mynamespace -l label=singlecontainer -f 100      

//Tail log from container named containername in pods with label muticontainer in ns mynamespace and follow the last 100 line logs 
./k8slog -n mynamespace -l label=muticontainer -f 100 -c containername    

//Tail log from container named containername in pods with label muticontainer in ns mynamespace and follow the recent 1.5 hour logs 
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername    

//Tail log from container named containername in pods with label muticontainer in ns mynamespace and follow the recent 1.5 hour logs and it is u
sed as pod inside k8s
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -r inside    

//Tail log from container named containername in pods with label muticontainer in ns mynamespace and follow the recent 1.5 hour logs with the ku
beconfig file
./k8slog -n mynamespace -l label=muticontainer -s 1.5h -c containername -kubeconfig admin.kubeconfig    
## Options:    
  -c string    
  &emsp; -c container name , Container name when multiple containers in pod    
  -f int    
  &emsp; -f 10,The number of lines from the end of the logs to show. Defaults to -1, showing all logs. (default -1)    
  -kubeconfig string    
  &emsp; (optional) absolute path to the kubeconfig file,default is ${HOME}/.kube/config    
  -l string    
  &emsp; -l app=nginx,Selector (label query) to filter on. If present, default to ".*" for the pod-query.    
  -n string    
  &emsp; -n namespace (default "default")    
  -r string    
  &emsp; -r inside means run inside k8s,by default it is running outside k8s (default "outside")    
  -s duration    
  &emsp; -s 10s : Return logs newer than a relative duration like 5s, 2m, or 2.5h.  
