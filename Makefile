.PHONY:	build binary

IMAGE = example-grpc
TAG = 0.0.1

build:
	docker build -t $(IMAGE):$(TAG) .

binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./greeter_server/

build-proto:
	protoc -I helloworld/ helloworld/*.proto --go_out=plugins=grpc:helloworld

clean:
	rm server
