package popcount

import (
	"fmt"
	"sync"
)

var (
	once sync.Once
	pc   [256]byte
)

func initTable() {
	fmt.Println("Initialization table")
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// Table returns initialized table for bit counting.
func Table() [256]byte {
	once.Do(initTable)
	return pc
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	table := Table()
	return int(pc[byte(x>>(0*8))] +
		table[byte(x>>(1*8))] +
		table[byte(x>>(2*8))] +
		table[byte(x>>(3*8))] +
		table[byte(x>>(4*8))] +
		table[byte(x>>(5*8))] +
		table[byte(x>>(6*8))] +
		table[byte(x>>(7*8))])
}
