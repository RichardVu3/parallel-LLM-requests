package datasets

import (
	"os"
	"bufio"
	"strconv"
)

const OUTPUTFILE = "datasets/outputs.txt"

type Input struct {
	ID int
	Prompt string
	Streaming bool
	Temperature float64
}

func NewInput(id int, prompt string, streaming bool, temperature float64) *Input {
	return &Input{
		ID: id,
		Prompt: prompt,
		Streaming: streaming,
		Temperature: temperature,
	}
}

type Reader struct {
	inputFile string
}

func NewReader(runMode string) *Reader {
	INPUTFILE := ""
	if runMode == "TEST" {
		INPUTFILE = "datasets/small-inputs.txt"
	} else {
		INPUTFILE = "datasets/inputs.txt"
	}
	return &Reader{inputFile: INPUTFILE}
}

func (r *Reader) GetInput() *[]Input {
	inputs := make([]Input, 0)
	file, err := os.Open(r.inputFile)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
        panic(err)
    }
	i := 1
    for scanner.Scan() {
		prompt := scanner.Text()
		inputs = append(inputs, *NewInput(i, prompt, false, 0))
		i++
    }
	return &inputs
}

type Output struct {
	ID int
	Response *string
}

func NewOutput(id int, response *string) *Output {
	return &Output{
		ID: id,
		Response: response,
	}
}

type Writer struct {
	outputFile string
}

func NewWriter() *Writer {
	return &Writer{outputFile: OUTPUTFILE}
}

func (w *Writer) Write(outputs *[]Output) {
	file, err := os.OpenFile(w.outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, output := range *outputs {
		if output.Response != nil {
			_, err = file.WriteString("ID " + strconv.Itoa(output.ID) + ": " + *output.Response + "\n\n")
			if err != nil {
				panic(err)
			}
		}
	}
}

func (w *Writer) Clear() {
	file, err := os.OpenFile(w.outputFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}