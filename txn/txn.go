package txn

import "sync"

type TransactionManager struct {
	next        int64
	active      map[int64]struct{}
	maxCommited int64
}

func NewTxnManager() *TransactionManager {
	txm := &TransactionManager{}
	txm.maxCommited = -1
	txm.active = make(map[int64]struct{})
	return txm
}

func (txm *TransactionManager) Next() *Transaction {
	active := make(map[int64]struct{})
	for key, _ := range txm.active {
		active[key] = struct{}{}
	}
	txn := newTransaction(txm.next, active, txm.maxCommited)
	txm.next++
	return txn
}

func (txm *TransactionManager) End(txn *Transaction) {
	id := txn.id
	if _, ok := txm.active[id]; ok {
		delete(txm.active, id)
	}
	if id > txm.maxCommited {
		txm.maxCommited = id
	}
}

type Transaction struct {
	mutex  sync.Mutex
	id     int64
	max    int64
	active map[int64]struct{}
}

func newTransaction(next int64, active map[int64]struct{}, maxCommited int64) *Transaction {
	return &Transaction{
		id:     next,
		max:    maxCommited,
		active: active,
	}
}

func (txn *Transaction) Id() int64 {
	return txn.id
}

func (txn *Transaction) Max() int64 {
	return txn.max
}

func (txn *Transaction) Active() map[int64]struct{} {
	return txn.active
}
