
WORKSPACE := ../../workspace
PATH := $(HOME)/go/bin:$(PATH)

.PHONY: all dep test clean
all: edraj.pb.go health.pb.go root $(WORKSPACE)/certs

${WORKSPACE}/certs: 
	certstrap --depot-path=${WORKSPACE}/certs init         --passphrase "" -o edraj -cn edrajRootCA
	certstrap --depot-path=${WORKSPACE}/certs request-cert --passphrase "" --domain root.edraj.io
	certstrap --depot-path=${WORKSPACE}/certs sign         --CA edrajRootCA root.edraj.io
	certstrap --depot-path=${WORKSPACE}/certs request-cert --passphrase "" --domain localhost
	certstrap --depot-path=${WORKSPACE}/certs sign         --CA edrajRootCA localhost
	certstrap --depot-path=${WORKSPACE}/certs request-cert --passphrase "" -cn admin
	certstrap --depot-path=${WORKSPACE}/certs sign         --CA edrajRootCA admin
	certstrap --depot-path=${WORKSPACE}/certs request-cert --passphrase "" -cn kefah
	certstrap --depot-path=${WORKSPACE}/certs sign         --CA edrajRootCA kefah

dep: 
	go get -u github.com/square/certstrap
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/hashicorp/consul
	go get -u github.com/hashicorp/consul/api
	go get -v -d ./...

edraj.pb.go: ../schema/edraj.proto 
	protoc -I ../schema/ ../schema/edraj.proto --go_out=plugins=grpc:. 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' ./edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' ./edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' ./edraj.pb.go

health.pb.go: ../schema/health.proto
	protoc -I ../schema/ ../schema/health.proto --go_out=plugins=grpc:. 
	sed -i 's/fileDescriptor0/fileDescriptor1/g' health.pb.go
	sed -i '1{s/$$/  DO NOT EDIT./}' ./health.pb.go


root: main.go edraj.pb.go
	go build -o root main.go edraj.pb.go health.pb.go

test: 
	go test

clean: 
	rm -f root edraj.pb.go
