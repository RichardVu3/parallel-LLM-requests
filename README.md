# Parallel Processing of Large Language Model Requests in Go

This project demonstrates the use of parallel programming techniques in Golang to process requests for a Large Language Model (LLM) efficiently. It focuses on optimizing performance for conversational AI by implementing two different parallelization approaches: **Simple Parallel** and **Work-Stealing Parallelism**.

## Features
- Written in **Go**, leveraging goroutines for efficient parallel processing.
- Uses the **LLama 3.2** model via the **Ollama framework**.
- Implements:
  - **Sequential Processing** for baseline comparison.
  - **Simple Parallel Processing** using pipelining and fixed thread allocation.
  - **Work-Stealing Parallelism** for dynamic load balancing.

## Directory Structure
```
parallel-llm-requests/
├── install/              # Dependency installation scripts
├── datasets/             # Input and output datasets
├── runner/               # Parallel and sequential execution modes
├── queue/                # Lock-free linked-list queue implementation
├── main.go               # Entry point
├── test/                 # Testing scripts
├── benchmark/            # Benchmarking scripts and analysis
```

### Key Files
- **`llm.go`**: Manages LLM interactions.
- **`datasets.go`**: Defines data structures for inputs and outputs.
- **`main.go`**: The entry point for running the program.

## Installation

### Prerequisites
- **Golang 1.23.3** or later.
- **LLama 3.2** model running in Ollama framework.

### Steps
1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/parallel-llm.git
   cd parallel-llm-requests/
   ```
2. Install dependencies:
   ```bash
   source install/go-install.sh
   source install/ollama-install.sh
   ```

## Usage
Run the program using the following commands:
```bash
go run parallel-llm-requests/main <runType> <threadNums>
```
- **`<runType>`**:
  - `seq`: Sequential processing.
  - `sp`: Simple parallel processing.
  - `ws`: Work-stealing parallelism.
- **`<threadNums>`**: Number of threads for parallel modes.

### Examples
- Sequential processing:
  ```bash
  go run parallel-llm-requests/main seq
  ```
- Simple parallel processing with 4 threads:
  ```bash
  go run parallel-llm-requests/main sp 4
  ```
- Work-stealing with 4 threads:
  ```bash
  go run parallel-llm-requests/main ws 4
  ```

## Testing
Use the test scripts to validate implementations:
```bash
bash test/test.sh <testType>
```
- **`<testType>`**:
  - `streaming`: Validates single prompt with streaming.
  - `nonstreaming`: Validates single prompt without streaming.
  - `workstealing`: Validates work-stealing correctness.

## Performance analysis
See [benchmark/performance_analysis.md](benchmark/performance_analysis.md).
