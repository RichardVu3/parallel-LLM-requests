#!/bin/bash

# Check if the user has provided the correct number of arguments
if [ "$#" -ne 1 ]; then
    echo "Usage: bash test/test.sh <mode>"
    exit 1
fi

# Check if the user has provided the correct mode
if [ "$1" != "streaming" ] && [ "$1" != "nonstreaming" ] && [ "$1" != "workstealing" ]; then
    echo "Invalid mode. Please use either streaming, nonstreaming or workstealing"
    exit 1
fi

# Activate golang1.23.3 and ollama

export PATH=$(pwd)/golang/go/bin:$(pwd)/ollama-package/bin:$PATH

if command -v ollama &>/dev/null; then
    echo "ollama already exists."
else
    echo "Please run source install/ollama-install.sh"
fi

ollama serve &

OLLAMA_PID=$! &&

echo "ollama is running in the background with PID: $OLLAMA_PID" &&

sleep 5 &&

ollama pull llama3.2 &&

wkdir=$(pwd) &&

# Go to the project directory
cd $wkdir &&

# Run the test.go file with the mode provided
go run $wkdir/test/test.go $1

# Kill the ollama server
kill -9 $OLLAMA_PID && echo "ollama stopped."