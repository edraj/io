
PATH := ~/go/bin:$(PATH)

.PHONY: all test dep clean help
all: edraj.pb.go server

# PKG_LIST := $(shell go list github.com/edraj/io/... | grep -v /vendor/)


edraj.pb.go: ../schema/edraj.proto 
	protoc -I ../schema/ ../schema/edraj.proto --go_out=plugins=grpc:. 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' ./edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' ./edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' ./edraj.pb.go

#./python_client/edraj_pb2_grpc.py: ../schema/edraj.proto
#	python3 -m grpc_tools.protoc --proto_path=../schema --python_out=./python_client/ --grpc_python_out=./python_client/ ../schema/edraj.proto

.PHONY: certs
certs: 
	certstrap --depot-path=../../workspace/certs init --passphrase "" -o edraj -cn edrajRootCA
	certstrap --depot-path=../../workspace/certs request-cert --passphrase "" --domain edraj.io
	certstrap --depot-path=../../workspace/certs sign --CA edrajRootCA edraj.io
	certstrap --depot-path=../../workspace/certs request-cert --passphrase "" --domain localhost
	certstrap --depot-path=../../workspace/certs sign --CA edrajRootCA localhost
	certstrap --depot-path=../../workspace/certs request-cert --passphrase "" -cn admin
	certstrap --depot-path=../../workspace/certs sign --CA edrajRootCA admin
	certstrap --depot-path=../../workspace/certs request-cert --passphrase "" -cn kefah
	certstrap --depot-path=../../workspace/certs sign --CA edrajRootCA kefah

server: main.go edraj.pb.go storage.go entryman.go grpc.go http.go
	go build -o server main.go edraj.pb.go storage.go entryman.go grpc.go http.go

dep: 
	go get gopkg.in/mgo.v2 
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/square/certstrap
	go get -v -d ./...

test: 
	go test

clean: 
	rm -f server edraj.pb.go

help:
	## Display this help screen 	
	#@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(PKG_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
