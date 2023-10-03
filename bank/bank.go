package bank

import (
	"bank/txn"
	"fmt"
	"time"
)

type Bank struct {
	oracle     *txn.TransactionManager
	accountNum int
	accounts   []*Account
}

func NewBank(accNum int, balance int64, manager *txn.TransactionManager) *Bank {
	accList := make([]*Account, 0, accNum)
	for i := 0; i < accNum; i++ {
		accList = append(accList, newAccount(manager, balance))
	}
	return &Bank{
		accountNum: accNum,
		accounts:   accList,
		oracle:     manager,
	}
}

func (b *Bank) GetAccountsNum() int {
	return b.accountNum
}

func (b *Bank) Transfer(from, to, amount int64) {
	if from == to {
		return
	}

	fromAcc := b.accounts[from]
	toAcc := b.accounts[to]

	var min, max *Account
	if from < to {
		min = fromAcc
		max = toAcc
	} else {
		min = toAcc
		max = fromAcc
	}

	txn := b.oracle.Next()

	min.mutex.Lock()
	defer min.mutex.Unlock()
	max.mutex.Lock()
	defer max.mutex.Unlock()

	if b.canUpdate(fromAcc, txn) && b.canUpdate(toAcc, txn) {
		toAcc.update(txn, toAcc.balance(txn)+amount)
		fromAcc.update(txn, fromAcc.balance(txn)-amount)
	}

	b.oracle.End(txn)
}

func (b *Bank) Holdings() int64 {
	txn := b.oracle.Next()
	total := int64(0)
	for _, account := range b.accounts {
		total += account.balance(txn)
	}
	return total
}

func (b *Bank) AccountsNum() int {
	return b.accountNum
}

func (b *Bank) LongRunningRead() {
	txn := b.oracle.Next()
	holdings := int64(0)
	for _, account := range b.accounts {
		holdings += account.balance(txn)
		time.Sleep(100 * time.Millisecond)
	}

	b.oracle.End(txn)
	fmt.Printf("Holding: %v, currentTxn: %v, maxTxn: %v \n", holdings, txn.Id(), txn.Max())
}

func (b *Bank) canUpdate(acc *Account, txn *txn.Transaction) bool {
	return acc.canBeUpdateBy(txn)
}