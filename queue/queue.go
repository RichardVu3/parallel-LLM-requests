package queue

import (
	"sync/atomic"
	"parallel-llm-requests/datasets"
)

type Node struct {
	input *datasets.Input
	next atomic.Pointer[Node]
}

func NewNode(input *datasets.Input) *Node {
	return &Node{input: input}
}

type LockFreeQueue struct {
	ID int
	head atomic.Pointer[Node]
	tail atomic.Pointer[Node]
}

func NewLockFreeQueue(id int) *LockFreeQueue {
	var head atomic.Pointer[Node]
	var tail atomic.Pointer[Node]
	node := Node{}
	head.Store(&node)
	tail.Store(&node)
	newQueue := &LockFreeQueue{ID: id}
	newQueue.head.Store(head.Load())
	newQueue.tail.Store(tail.Load())
	return newQueue
}

func (queue *LockFreeQueue) Enqueue(task *datasets.Input) {
	newNode := NewNode(task)
	for {
		tail := queue.tail.Load()
		next := tail.next.Load()
		if tail == queue.tail.Load() {
			if next == nil {
				if tail.next.CompareAndSwap(next, newNode) {
					queue.tail.CompareAndSwap(tail, newNode)
					return
				}
			} else {
				queue.tail.CompareAndSwap(tail, next)
			}
		}
	}
}

func (queue *LockFreeQueue) Dequeue() *datasets.Input {
	var req *datasets.Input
	for {
		head := queue.head.Load()
		tail := queue.tail.Load()
		first := head.next.Load()
		if head == queue.head.Load() {
			if head == tail {
				if first == nil {
					return req
				}
				queue.tail.CompareAndSwap(tail, first)
			} else {
				request := first.input
				if queue.head.CompareAndSwap(head, first) {
					return request
				}
			}
		}
	}
}

func (queue *LockFreeQueue) IsEmpty() bool {
	head := queue.head.Load()
	tail := queue.tail.Load()
	return head == tail && head.next.Load() == nil
}