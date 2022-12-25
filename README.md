# homework3
### поднять миникуб
minikube start driver=docker

### поключить расширения игрес
minikube addons enable ingress
minikube addons enable ingress-dns

### прокинуть порты
minikube tunnel

### установить istio
istioctl install

### 
istioctl operator init --watchedNamespaces istio-system --operatorNamespace istio-operator
kubectl apply -f ./k8s/apigateway/istio.yaml

### установить metallb
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml

### настроить metallb
kubectl apply -f ./k8s/apigateway/metallb.yaml

### создать неймспейс приложения, настроить его для использования с istio и установить приложение
kubectl create namespace userservice

kubectl label namespace userservice istio-injection=enabled --overwrite

helm install --set name=userservice userservice --namespace userservice ./k8s/userservice

### установить apigateway
kubectl apply -f ./k8s/apigateway/gateway.yaml

### посмотреть <EXTERNAL-IP> и прописать его в /etc/hosts
kubectl get svc istio-ingressgateway -n istio-system

sudo nano /etc/hosts