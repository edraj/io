#!/bin/bash
sudo dnf install git-lfs protobuf-compiler mongodb mongodb-server golang mongo-tools
sudo systemctl start mongod
git lfs pull

pip3 install --user grpcio-tools
export GOPATH="$HOME/go"
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get gopkg.in/mgo.v2 

mongo edraj --eval 'db.content.drop();'
# mongodump -d edraj --gzip
mongorestore -d edraj --gzip --dir ../sampledata/edraj
mongo edraj --eval 'db.content.count();'
mongo edraj --eval 'db.content.getIndices();'
#mongo edraj --eval 'db.content.createIndex( { displayname: "text", body: "text" }, {"default_language":"none"} );' #  Arabic sadly requires 3rd party license
