package main

import (
	"os"
	"fmt"
	"strconv"
	"parallel-llm-requests/runner"
)

const RUNMODE = "TEST" // TEST means write output to file

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main <runType> <threadNums>")
		os.Exit(1)
	}
	threadNums := 0
	if len(os.Args) == 3 {
		threadNums, _ = strconv.Atoi(os.Args[2])
	}
	processor := runner.NewRunner(os.Args[1], threadNums, RUNMODE)
	processor.Run()
}
