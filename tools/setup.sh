#!/bin/bash
sudo dnf install git-lfs protobuf-compiler mongodb mongodb-server golang
sudo systemctl start mongodb

pip3 install --user grpcio-tools
export GOPATH="$HOME/go/bin:$PATH"
go get github.com/BurntSushi/toml gopkg.in/mgo.v2 github.com/gorilla/mux
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/golang/protobuf/protoc-gen-go
mongo edraj --eval 'db.content.drop();'
# mongodump -d edraj --gzip
mongorestore -d edraj --gzip --dir ../sampledata/
mongo edraj --eval 'db.content.count();'
mongo edraj --eval 'db.content.getIndices();'
#mongo edraj --eval 'db.content.createIndex( { displayname: "text", body: "text" }, {"default_language":"none"} );' #  Arabic sadly requires 3rd party license
