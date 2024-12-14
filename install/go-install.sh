#!/bin/bash

wkdir=$(pwd)

mkdir -p $wkdir/golang &&

curl -L https://go.dev/dl/go1.23.3.linux-amd64.tar.gz -o $wkdir/golang/go1.23.3.linux-amd64.tar.gz &&

echo "Downloaded successfully. Unzipping ..." &&

tar -C $wkdir/golang -xzf $wkdir/golang/go1.23.3.linux-amd64.tar.gz &&

export PATH=$wkdir/golang/go/bin:$PATH &&

go get "github.com/tmc/langchaingo/llms" &&

echo "Go 1.23.3 installed successfully."