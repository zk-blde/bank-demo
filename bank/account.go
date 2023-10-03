package bank

import (
	"bank/txn"
	"sync"
)

type Account struct {
	mutex   sync.Mutex
	history []*AccountEntry
	lastTxn *txn.Transaction
}

func newAccount(txm *txn.TransactionManager, balance int64) *Account {
	tx := txm.Next()
	acc := &Account{}
	acc.update(tx, balance)
	txm.End(tx)
	return acc
}

func (acc *Account) update(txn *txn.Transaction, amount int64) {
	acc.lastTxn = txn
	acc.history = append(acc.history, newAccountEntry(txn, amount))
}

func (acc *Account) balance(txn *txn.Transaction) int64 {
	tid := txn.Id()
	maxid := txn.Max()
	active := txn.Active()

	entries := acc.history
	var res *AccountEntry
	for _, entry := range entries {
		if entry.txn.Id() >= tid {
			continue
		}

		if entry.txn.Id() > maxid {
			continue
		}

		if _, ok := active[entry.txn.Id()]; ok {
			continue
		}

		res = entry

	}
	return res.balance
}

func (acc *Account) canBeUpdateBy(txn *txn.Transaction) bool {
	if acc.lastTxn.Id() > txn.Id() {
		return false
	}
	active := txn.Active()

	for _, entry := range acc.history {
		if _, ok := active[entry.txn.Id()]; ok {
			return false
		}
	}
	return true
}

type AccountEntry struct {
	txn     *txn.Transaction
	balance int64
}

func newAccountEntry(txn *txn.Transaction, amount int64) *AccountEntry {
	return &AccountEntry{
		txn:     txn,
		balance: amount,
	}
}
