package popcount

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func popCount1(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func popCount2(x uint64) int {
	num := 0
	for i := uint8(0); i < 64; i++ {
		if x&(1<<i) != 0 {
			num++
		}
	}
	return num
}

func popCount3(x uint64) int {
	num := 0
	for x != 0 {
		x = x & (x - 1)
		num++
	}
	return num
}

func popCount4(x uint64) int {
	var num byte
	var i uint
	for i = 0; i < 8; i++ {
		num += pc[byte(x>>(i*8))]
	}
	return int(num)
}
