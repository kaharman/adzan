export GO111MODULE=on
NAME=adzan-service
TAG := $(shell git describe --candidates=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
VERSION := $(TAG:v%=%)
ifeq ($(TAG),)
	VERSION := $(BRANCH)
endif

build:
	protoc --plugin=protoc-gen-doc=${HOME}/go/bin/protoc-gen-doc --doc_out=adzan/delivery/http/grpcdoc --doc_opt=adzan/delivery/http/grpcdoc/html.tmpl,index.html --proto_path=adzan/delivery/http/grpcdoc grpc.proto
	go mod tidy
	go build -o ${NAME} app/*.go

clean:
	if [ -f ${NAME} ] ; then rm ${NAME} ; fi

docker:
	docker build -t ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION) .

run:
	protoc --plugin=protoc-gen-doc=${HOME}/go/bin/protoc-gen-doc --doc_out=adzan/delivery/http/grpcdoc --doc_opt=adzan/delivery/http/grpcdoc/html.tmpl,index.html --proto_path=adzan/delivery/http/grpcdoc grpc.proto
	go mod tidy
	go run app/*.go -c config.yaml

push:
	docker push ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION)
