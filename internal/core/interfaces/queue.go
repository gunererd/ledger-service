package interfaces

import "ledger-service/internal/core/types"

type Queue interface {
	Enqueue(t types.Transaction)
	Dequeue() types.Transaction
}
