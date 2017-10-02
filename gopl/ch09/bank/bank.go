package bank

type Withdrawal struct {
	amount  int
	succeed chan bool
}

var (
	deposits    = make(chan int)
	balances    = make(chan int)
	withdrawals = make(chan Withdrawal)
)

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func Withdraw(amount int) bool {
	ch := make(chan bool)
	withdrawals <- Withdrawal{amount, ch}
	return <-ch
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case w := <-withdrawals:
			if w.amount > balance || w.amount <= 0 {
				w.succeed <- false
				continue
			}
			balance -= w.amount
			w.succeed <- true
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}
