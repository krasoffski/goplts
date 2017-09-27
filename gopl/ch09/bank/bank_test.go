package bank

import (
	"sync"
	"testing"
)

func TestBank(t *testing.T) {
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(amount int) {
			Deposit(amount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if got, want := Balance(), (1000+1)*1000/2; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
