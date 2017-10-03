package bank

import (
	"sync"
	"testing"
)

func TestDeposit(t *testing.T) {
	balance := Balance()
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(amount int) {
			Deposit(amount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if got, want := Balance(), (1000+1)*1000/2+balance; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

func TestWithdraw(t *testing.T) {
	balance := Balance()
	Deposit((1000 + 1) * 1000 / 2)
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(amount int) {
			Withdraw(amount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if got, want := Balance(), balance; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

func TestWithdrawNegative(t *testing.T) {
	deposit := 11111
	balance := Balance()
	Deposit(deposit)

	if result := Withdraw(-17); result != false {
		t.Errorf("Result = %v, want %v", result, false)
	}

	if got, want := Balance(), balance+deposit; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
