
PATH := $(HOME)/go/bin:$(PATH)

.PHONY: all test dep clean proto
all: go/edraj.pb.go go/client python/edraj_pb2_grpc.py node/edraj_grpc_pb.js

proto: go/edraj.pb.go python/edraj_pb2_grpc.py node/edraj_grpc_pb.js

go/edraj.pb.go: ../schema/edraj.proto 
	protoc -I ../schema/ ../schema/edraj.proto --go_out=plugins=grpc:go 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' go/edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' go/edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' go/edraj.pb.go

python/edraj_pb2_grpc.py: ../schema/edraj.proto
	python3 -m grpc_tools.protoc --proto_path=../schema --python_out=./python/ --grpc_python_out=./python/ ../schema/edraj.proto

node/edraj_grpc_pb.js: ../schema/edraj.proto
	./node/node_modules/.bin/grpc_tools_node_protoc --js_out=import_style=commonjs,binary:./node/ --grpc_out=./node/ --plugin=protoc-gen-grpc=./node/node_modules/grpc-tools/bin/grpc_node_plugin --proto_path=../schema ../schema/edraj.proto

go/client: go/main.go go/edraj.pb.go
	go build -o go/client go/main.go go/edraj.pb.go 

dep: 
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -v -d ./...
	pip3 install --user grpcio-tools
	cd node && npm install ; cd ..

test: 
	go test

clean: 
	rm -f  go/edraj.pb.go python/edraj_pb2_grpc.py python/edraj_pb2.py node/edraj*pb.js
