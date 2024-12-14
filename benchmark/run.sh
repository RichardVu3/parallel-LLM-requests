#!/bin/bash

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

mkdir -p $wkdir/benchmark/runtime &&

touch $wkdir/benchmark/runtime/sequential.txt &&

echo "Start sequential run..." &&

for i in {1..5}
do
    go run main seq >> $wkdir/benchmark/runtime/sequential.txt
done &&

mkdir -p $wkdir/benchmark/runtime/simple-parallel &&

echo "Start simple parallel run..." &&

for thread in 2 4 6 8 12
do
    touch $wkdir/benchmark/runtime/simple-parallel/$thread.txt
    for i in {1..5}
    do
        go run main sp $thread >> $wkdir/benchmark/runtime/simple-parallel/$thread.txt
    done
done &&

mkdir -p $wkdir/benchmark/runtime/work-stealing &&

echo "Start work stealing run..." &&

for thread in 2 4 6 8 12
do
    touch $wkdir/benchmark/runtime/work-stealing/$thread.txt
    for i in {1..5}
    do
        go run main ws $thread >> $wkdir/benchmark/runtime/work-stealing/$thread.txt
    done
done &&

kill -9 $OLLAMA_PID && echo "ollama stopped."

python3 $wkdir/benchmark/analysis.py

echo "Finished."