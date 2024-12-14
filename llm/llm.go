package llm

import (
	"context"
	"fmt"
	"log"
	"parallel-llm-requests/datasets"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const MODEL = "llama3.2"

type Agent struct {
	llm *ollama.LLM
	ctx *context.Context
}

func NewAgent() *Agent {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	return &Agent{llm: llm, ctx: &ctx}
}

func (a *Agent) Invoke(input *datasets.Input) (output *datasets.Output) {
	output = &datasets.Output{ID: input.ID}
	errorResponse := "empty response due to an error"
	emptyResponse := ""
	if input.Streaming {
		_, err := llms.GenerateFromSinglePrompt(
			*a.ctx, a.llm, input.Prompt,
			llms.WithTemperature(input.Temperature),
			llms.WithStreamingFunc(
				func(ctx context.Context, chunk []byte) error {
					fmt.Print(string(chunk))
					return nil
				},
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		output.Response = &emptyResponse
	} else {
		completion, err := llms.GenerateFromSinglePrompt(
			*a.ctx, a.llm, input.Prompt, llms.WithTemperature(input.Temperature),
		)
		if err != nil {
			log.Fatal(err)
			output.Response = &errorResponse
		} else {
			output.Response = &completion
		}
	}
	return output
}
