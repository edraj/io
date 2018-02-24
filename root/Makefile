EDRAJ_WORKSPACE := ../../workspace
# PKG := "gitlab.com/edraj/io"
# PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all dep test clean
all: edraj.pb.go root $(EDRAJ_WORKSPACE)/certs

${EDRAJ_WORKSPACE}/certs: 
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs init         --passphrase "" -o edraj -cn edrajRootCA
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs request-cert --passphrase "" --domain root.edraj.io
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs sign         --CA edrajRootCA root.edraj.io
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs request-cert --passphrase "" --domain localhost
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs sign         --CA edrajRootCA localhost
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs request-cert --passphrase "" -cn admin
	certstrap --depot-path=${EDRAJ_WORKSPACE}/certs sign          --CA edrajRootCA admin

dep: 
	go get -u github.com/square/certstrap
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -v -d ./...

edraj.pb.go: ../schema/edraj.proto 
	protoc -I ../schema/ ../schema/edraj.proto --go_out=plugins=grpc:. 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' ./edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' ./edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' ./edraj.pb.go


root: main.go edraj.pb.go
	go build -o root main.go edraj.pb.go

test: 
	go test

clean: 
	rm -f root edraj.pb.go
