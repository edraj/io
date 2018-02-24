
.PHONY: all test dep clean proto
all: go/edraj.pb.go go/client python/edraj_pb2_grpc.py

proto: go/edraj.pb.go python/edraj_pb2_grpc.py

go/edraj.pb.go: ../schema/edraj.proto 
	protoc -I ../schema/ ../schema/edraj.proto --go_out=plugins=grpc:go 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' go/edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' go/edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' go/edraj.pb.go

python/edraj_pb2_grpc.py: ../schema/edraj.proto
	python3 -m grpc_tools.protoc --proto_path=../schema --python_out=./python/ --grpc_python_out=./python/ ../schema/edraj.proto


go/client: go/client.go go/edraj.pb.go
	go build -o go/client go/client.go go/edraj.pb.go 

dep: 
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -v -d ./...
	pip3 install --user grpcio-tools

test: 
	go test

clean: 
	rm -f  go/edraj.pb.go python/edraj_pb2_grpc.py python/edraj_pb2.py
