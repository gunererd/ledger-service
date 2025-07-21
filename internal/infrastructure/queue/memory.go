package queue

import "ledger-service/internal/core/types"

type InMemoryQueue struct {
	ch chan types.Transaction
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		ch: make(chan types.Transaction, 100), // Buffered channel
	}
}

func (q *InMemoryQueue) Enqueue(t types.Transaction) {
	q.ch <- t
}

func (q *InMemoryQueue) Dequeue() types.Transaction {
	return <-q.ch // Blocks until item available
}
