#!/bin/sh
export GO111MODULE=on
GITHUB_USER="hellgate75"
PROJECT_NAME="go-deploy"
BASE_FOLDER="$GOPATH/src/github.com/$GITHUB_USER"
PROJECT_FOLDER="$BASE_FOLDER/$PROJECT_NAME"
echo "Creating base folder $BASE_FOLDER into folder: GOPATH $GOPATH"
mkdir -p $BASE_FOLDER
echo "Linking project into folder: GOPATH $GOPATH"
ln -s $PROJECT_FOLDER ./
echo "Changing folder to $PROJECT_FOLDER"
cd $PROJECT_FOLDER
echo "Running go procedure into folder:$PROJECT_FOLDER"
go mod init
go mod tidy
go test -v ./...
go build -v ./...