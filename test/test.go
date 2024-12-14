package main

import (
	"fmt"
	"os"
	"parallel-llm-requests/datasets"
	"parallel-llm-requests/llm"
	"parallel-llm-requests/runner"
)

const RUNMODE = "TEST"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test/test.go <mode>")
		fmt.Println("mode can be: streaming, nonstreaming, or workstealing")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "streaming":
		fmt.Println("Prompt: Give me one sentence about the University of Chicago")
		agent := llm.NewAgent()
		input := datasets.NewInput(1, "Give me one sentence about the University of Chicago", true, 0.5)
		agent.Invoke(input)

	case "nonstreaming":
		fmt.Println("Prompt: Give me one sentence about the University of Chicago")
		agent := llm.NewAgent()
		input := datasets.NewInput(1, "Give me one sentence about the University of Chicago", false, 0.5)
		output := agent.Invoke(input)
		fmt.Println("Output:", *output.Response)

	case "workstealing":
		fmt.Println("Running workstealing with 5 threads")
		threadNums := 5
		runType := "ws"
		processor := runner.NewRunner(runType, threadNums, RUNMODE)
		processor.Run()
		fmt.Println("Process finished successfully. Please check datasets/outputs.txt for output.")

	default:
		fmt.Println("Unknown mode:", mode, ". Please use streaming, nonstreaming, or workstealing.")
	}
}
