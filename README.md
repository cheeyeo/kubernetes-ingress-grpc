# Kubernetes Ingress with GRPC

This is an example repo on how to integrate (Ingress Controller)[https://github.com/kubernetes/ingress-nginx] with minikube on localhost and running a secure grpc service with it.

This is to test the latest nginx functionality to allow for grpc proxying:
https://www.nginx.com/blog/nginx-1-13-10-grpc/

There are 2 main approaches to allowing for GRPCs: terminating at nginx level; or proxying the request to the app. This repo opted for the second option as I could not get the first option to work.

# Enabling ingress addon in minikube

minikube comes with an ingress controller which you can enable automatically.

This will not work as it is not running the latest version of nginx to test the grpc proxing.

I have create a folder `nginx-ingress-minikube` which contains the same addon files from minnikube but with the `nginx-ingress-controller` pointing to `quay.io/aledbf/nginx-ingress-controller:0.353`

To set it up, run the following:
```
minikube disable addon ingress

kubectl create namespace ingress-nginx

kubectl create -f nginx-ingress-minikube/ingress-configmap.yaml

kubectl create -f nginx-ingress-minikube/ingress-default-backend.yaml

kubectl create -f nginx-ingress-minikube/ingress-rbac.yaml

kubectl create -f nginx-ingress-minikube/ingress-rc.yaml
```

Check that all pods are up and running:
```
kubectl get pods -n ingress-nginx
```

If all is well you should see the following:
```
NAME                                     READY     STATUS    RESTARTS   AGE
default-http-backend-55c6c69b88-cbwr6    1/1       Running   0          4h
nginx-ingress-controller-79f5774-rsbjv   1/1       Running   1          4h

```
## Building the application

The example grpc app is the `helloworld` example from the grpc repo itself. It echoes back whatever message you pass to it. I have made some changes to the server to enable TLS and added some unary interceptors for logging purposes.

To build the server:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./greeter_server/
```

The example expects tls crt and key within the certs directory. When creating the certs, under `Common Name` make sure that you use a wildcard domain such as `*.example.com` and update the `/etc/hosts` file accordingly :
```
mkdir -p certs

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/tls.key -out certs/tls.crt
```

To publish it into minikube docker registry:
```
(https://github.com/kubernetes/minikube/blob/master/docs/reusing_the_docker_daemon.md)

eval $(minikube docker-env)

docker build -t example-grpc:0.0.1 -f Dockerfile .
```

To setup the service and ingress on the cluster:
```
kubectl create namespace grpc

kubectl create secret tls tls-secret --key certs/tls.key --cert certs/tls.crt -n grpc

kubectl create -f greeter-svc.yaml -n grpc

kubectl create -f greeter-ingress.yaml -n grpc

kubectl get pods -n grpc

NAME                              READY     STATUS    RESTARTS   AGE
greeter-server-85dbdfbf77-wn2ts   1/1       Running   0          20m
```

To run the client, update the `address` field accordingly:

```
go run greeter_client/main.go "This is a test"
```

If its working, it should echo back the response:
```
2018/04/11 17:14:14 Greeting: Hello Hey it works
2018/04/11 17:14:14 Greeting: Hello again Hey it works
```

Check the logs on the pod:
```
kubectl logs -f <id of greeter-server pod> -n grpc
```

Note that since we are using `nginx.ingress.kubernetes.io/ssl-passthrough` for the ingress, the logging is not showing up in the nginx-controller. This is a known issue as of this writing:

https://github.com/kubernetes/ingress-nginx/issues/2329


# TODO:

* Makefile to automate some of the build process above

* Use grpcurl to test for grpc services endpoints

* Extend example to include multiple grpc services

* Figure out why ssl termination does not work


# References:

(Nginx grpc support announce)[https://www.nginx.com/blog/nginx-1-13-10-grpc/]

(Ingress nginx grpc PR)[https://github.com/kubernetes/ingress-nginx/pull/2307]
