#!/bin/sh
export GO111MODULE=on
GITHUB_USER="hellgate75"
PROJECT_NAME="go-deploy"
BUILD_MODE="exe"
EXTENSION=""
BASE_FOLDER="$GOPATH/src/github.com/$GITHUB_USER"
PROJECT_FOLDER="$BASE_FOLDER/$PROJECT_NAME"
echo "Working dir: $(pwd)"
echo "Creating base folder '$PROJECT_FOLDER' into folder: GOPATH '$GOPATH'"
mkdir -p $PROJECT_FOLDER
echo "Linking project folder: '$PROJECT_FOLDER' into folder: GOPATH '$GOPATH'"
ln -s $(pwd) $PROJECT_FOLDER
echo "Changing folder to $PROJECT_FOLDER"
cd $PROJECT_FOLDER
echo "Running go procedure into folder:$PROJECT_FOLDER"
go mod init
go mod tidy
echo "Testing project into folder:$PROJECT_FOLDER"
go test -v .
OUT_FILE_NAME="$PROJECT_NAME-$(uname -o)-$(uname -m)$EXTENSION"
echo "Building project for making: $OUT_FILE_NAME"
go build -v -buildmode "$BUILD_MODE" -o "$OUT_FILE_NAME"