install-kubectl:
	curl -LO "https://dl.k8s.io/release/v1.32.0/bin/linux/amd64/kubectl"
	sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl

minikube-install:
	curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
	sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64

install-go:
	curl -LO https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
	rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
	rm go1.23.4.linux-amd64.tar.gz

# requires few minutes
start:
	minikube start

build-push-minikube:
	eval $(minikube docker-env)
	cd src && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o demo . && mv demo ../
	docker build . -t super_custom_demo_1337

expose:
	kubectl port-forward service/demo-service 8080:8080