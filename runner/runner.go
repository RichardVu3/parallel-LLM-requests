package runner

import (
	"fmt"
	"math/rand"
	"parallel-llm-requests/datasets"
	"parallel-llm-requests/llm"
	"parallel-llm-requests/queue"
	"slices"
	"sync"
	"time"
)

type Runner struct {
	runType    string
	threadNums int
	reader     *datasets.Reader
	writer     *datasets.Writer
	runMode    string
}

func NewRunner(runType string, threadNums int, runMode string) *Runner {
	runner := &Runner{
		runType:    runType,
		threadNums: threadNums,
		reader:     datasets.NewReader(runMode),
		writer:     datasets.NewWriter(),
		runMode:    runMode,
	}
	runner.writer.Clear()
	return runner
}

func (r *Runner) Run() {
	start := time.Now()
	switch r.runType {
	case "seq": // sequential
		r.runSequential()
	case "sp": // simple parallel
		r.runSimParallel()
	case "ws": // worksteal
		r.runWorksteal()
	default:
		panic("Invalid run type. Currently supported: seq, sp, ws")
	}
	fmt.Println(time.Since(start).Seconds())
}

func (r *Runner) runSequential() {
	inputs := r.reader.GetInput()
	agent := llm.NewAgent()
	outputs := make([]datasets.Output, len(*inputs))
	for i := range *inputs {
		agent.Invoke(&(*inputs)[i])
		output := agent.Invoke(&(*inputs)[i])
		outputs[i] = *output
	}
	if r.runMode == "TEST" {
		r.writer.Write(&outputs)
	}
}

func (r *Runner) runSimParallel() {
	inputs := r.reader.GetInput()

	threadNums := r.threadNums
	taskIndexes := make([][]int, threadNums)
	for i := 0; i < len(*inputs); i++ {
		taskIndexes[i%threadNums] = append(taskIndexes[i%threadNums], i)
	}

	generator := func(done <-chan struct{}, taskIndexes ...[]int) <-chan []int {
		intStream := make(chan []int)
		go func() {
			defer close(intStream)
			for _, taskIndex := range taskIndexes {
				select {
				case <-done:
					return
				case intStream <- taskIndex:
				}
			}
		}()
		return intStream
	}

	invoke := func(
		done <-chan struct{},
		intStream <-chan []int,
		threadNums int,
	) <-chan *[]datasets.Output {
		agent := llm.NewAgent()
		outputStream := make(chan *[]datasets.Output)

		var wg sync.WaitGroup

		worker := func(taskIndex []int) {
			defer wg.Done()
			outputs := make([]datasets.Output, len(taskIndex))
			for i := range taskIndex {
				output := agent.Invoke(&(*inputs)[taskIndex[i]])
				outputs[i] = *output
			}
			select {
			case <-done:
				return
			case outputStream <- &outputs:
			}
		}

		wg.Add(threadNums)
		for i := 0; i < threadNums; i++ {
			go func() {
				for taskIndex := range intStream {
					worker(taskIndex)
				}
			}()
		}

		go func() {
			wg.Wait()
			close(outputStream)
		}()

		return outputStream
	}

	done := make(chan struct{})
	defer close(done)

	intStream := generator(done, taskIndexes...)
	pipeline := invoke(done, intStream, threadNums)

	for output := range pipeline {
		if r.runMode == "TEST" {
			r.writer.Write(output)
		}
	}
}

func (r *Runner) runWorksteal() {
	inputs := r.reader.GetInput()
	threadNums := r.threadNums
	taskIndexes := make([]*queue.LockFreeQueue, threadNums)
	var wg sync.WaitGroup

	for i := 0; i < threadNums; i++ {
		taskIndexes[i] = queue.NewLockFreeQueue(i)
	}
	for i := 0; i < len(*inputs); i++ {
		taskIndexes[i%threadNums].Enqueue(&(*inputs)[i])
	}
	allOutputs := make([]*[]datasets.Output, threadNums)
	wg.Add(threadNums)
	for i := 0; i < threadNums; i++ {
		go func(ID int, taskIndexes *[]*queue.LockFreeQueue, allOutputs *[]*[]datasets.Output) {
			defer wg.Done()
			agent := llm.NewAgent()
			outputs := make([]datasets.Output, 0)
			// Finish its tasks first
			for !(*taskIndexes)[ID].IsEmpty() {
				input := (*taskIndexes)[ID].Dequeue()
				output := agent.Invoke(input)
				outputs = append(outputs, *output)
			}
			if r.runMode == "TEST" {
				fmt.Println("Thread", ID, "finished its tasks")
			}
			// Steal tasks from other threads
			emptyIDs := make([]int, 0)
			for {
				if len(emptyIDs) == threadNums-1 {
					break
				}
				stealID := rand.Intn(threadNums)
				if stealID == ID || slices.Contains(emptyIDs, stealID) {
					continue
				}
				if (*taskIndexes)[stealID].IsEmpty() {
					emptyIDs = append(emptyIDs, stealID)
					continue
				}
				input := (*taskIndexes)[stealID].Dequeue()
				if r.runMode == "TEST" {
					fmt.Println("Thread", ID, "steals task", input.ID, "from thread", stealID)
				}
				output := agent.Invoke(input)
				outputs = append(outputs, *output)
			}
			(*allOutputs)[ID] = &outputs
		}(i, &taskIndexes, &allOutputs)
	}
	wg.Wait()
	if r.runMode == "TEST" {
		for i := 0; i < threadNums; i++ {
			r.writer.Write(allOutputs[i])
		}
	}
}
