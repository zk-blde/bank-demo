package main

import (
	"bank/bank"
	"bank/txn"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func (demo *BankDemo) doLongRunningRead() {
	go func() {
		fmt.Println("Long Reader is running")
		for i := 0; i < 1000; i++ {
			demo.bank.LongRunningRead()

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (demo *BankDemo) printHoldings() {
	go func() {
		fmt.Println("printer is running")
		for i := 0; i < 1000; i++ {
			fmt.Println("")
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func (demo *BankDemo) doRandomTransfer() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
			for j := 0; j < 10000; j++ {
				demo.doRealRandomTransfer()
			}
			time.Sleep(5 * time.Millisecond)
			<-ch
		}(i)
	}
	wg.Wait()
}

func (demo *BankDemo) doRealRandomTransfer() {
	accNum := demo.bank.AccountsNum()
	from := rand.Intn(int(accNum))
	to := rand.Intn(int(accNum))

	amount := rand.Intn(100)
	demo.bank.Transfer(int64(from), int64(to), int64(amount))
}

type BankDemo struct {
	bank *bank.Bank
}

func newDemo(accNum int, balance int64) *BankDemo {
	return &BankDemo{bank: bank.NewBank(accNum, balance, txn.NewTxnManager())}
}

func main() {
	var wg sync.WaitGroup
	numAccount, balance := 100, int64(20)
	demo := newDemo(numAccount, balance)

	demo.printHoldings()
	wg.Add(1)
	demo.doLongRunningRead()
	wg.Add(1)
	demo.doRandomTransfer()
	wg.Wait()
}
