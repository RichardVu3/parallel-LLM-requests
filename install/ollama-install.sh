#!/bin/bash

wkdir=$(pwd)/ollama-package

mkdir -p $wkdir &&

curl -L https://ollama.com/download/ollama-linux-amd64.tgz -o $wkdir/ollama-linux-amd64.tgz &&

echo "Downloaded successfully. Unzipping ..." &&

tar -C $wkdir -xzf $wkdir/ollama-linux-amd64.tgz &&

export PATH=$PATH:$wkdir/bin &&

ollama serve &

OLLAMA_PID=$! &&

echo "ollama is running in the background with PID: $OLLAMA_PID" &&

sleep 5 &&

ollama pull llama3.2 &&

kill -9 $OLLAMA_PID && echo "ollama stopped."

echo "Ollama installed successfully."
