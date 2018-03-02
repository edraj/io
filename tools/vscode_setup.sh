#!/bin/bash
go get -v github.com/uudashr/gopkgs/cmd/gopkgs
go get -v sourcegraph.com/sqs/goreturns
go get -u github.com/derekparker/delve/cmd/dlv

go get -v github.com/nsf/gocode
go get -v github.com/ramya-rao-a/go-outline
go get -v github.com/acroca/go-symbols
go get -v golang.org/x/tools/cmd/guru
go get -v golang.org/x/tools/cmd/gorename
go get -v github.com/rogpeppe/godef
go get -v golang.org/x/tools/cmd/godoc
go get -v github.com/golang/lint/golint

#sudo sh -c 'echo -e "[code]\nname=Visual Studio Code\nbaseurl=https://packages.microsoft.com/yumrepos/vscode\nenabled=1\ngpgcheck=1\ngpgkey=https://packages.microsoft.com/keys/microsoft.asc" > /etc/yum.repos.d/vscode.repo'
#sudo dnf install code
