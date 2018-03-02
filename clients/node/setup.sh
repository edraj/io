#!/bin/bash

sudo dnf install nodejs

npm init -y
npm install --save request body-parser lowdb express
npm install --save grpc grpc-tools grpc-health-check

node index.js

./node_modules/.bin/grpc_tools_node_protoc --js_out=import_style=commonjs,binary:./ --grpc_out=./ --plugin=protoc-gen-grpc=./node_modules/grpc-tools/bin/grpc_node_plugin --proto_path=../../schema ../../schema/edraj.proto

#grpc_tools_node_protoc --js_out=import_style=commonjs,binary:./ --grpc_out=./ --plugin=protoc-gen-grpc=./node_modules/grpc-tools/bin/grpc_node_plugin ../../schema/edraj.proto
#./node_modules/grpc-tools/bin/protoc
#grpc_tools_node_protoc --js_out=import_style=commonjs,binary:../node/static_codegen/route_guide/ --grpc_out=../node/static_codegen/route_guide/ --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` route_guide.proto

