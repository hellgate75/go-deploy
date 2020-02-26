#!/bin/sh
MODULE_PATH="./mod"
if [ "" != "$1" ]; then
	MODULE_PATH="$1"
fi
OS="$(sh ./os.sh)"
echo "OS: $OS"

if [ ! -e $GOPATH/src/github.com/hellgate75/go-deploy ]; then
	go get github.com/hellgate75/go-deploy
fi

if [ "windows" = "$OS" ]; then
	mkdir -p $MODULE_PATH/shell
	go build -v -o $MODULE_PATH/shell/shell.dll -buildmode=c-shared github.com/hellgate75/go-deploy/modules/shell
	mkdir -p $MODULE_PATH/service
	go build -v -o $MODULE_PATH/service/service.dll -buildmode=c-shared github.com/hellgate75/go-deploy/modules/shell
	mkdir -p $MODULE_PATH/copy
	go build -v -o $MODULE_PATH/copy/copy.dll -buildmode=c-shared github.com/hellgate75/go-deploy/modules/shell
	
	go build -o ./go-deploy.exe github.com/hellgate75/go-deploy
else
	mkdir -p $MODULE_PATH/shell
	go build -v -o $MODULE_PATH/shell/shell.so -buildmode=plugin github.com/hellgate75/go-deploy/modules/shell
	mkdir -p $MODULE_PATH/service
	go build -v -o $MODULE_PATH/service/service.so -buildmode=plugin github.com/hellgate75/go-deploy/modules/shell
	mkdir -p $MODULE_PATH/copy
	go build -v -o $MODULE_PATH/copy/copy.so -buildmode=plugin github.com/hellgate75/go-deploy/modules/shell
	
	go build -o ./go-deploy github.com/hellgate75/go-deploy

fi
