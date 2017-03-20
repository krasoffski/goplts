package intset

import "testing"

func TestLen(t *testing.T) {
	var x IntSet
	var length int
	x.Add(1)
	x.Add(144)
	x.Add(9)
	length = x.Len()
	if length != 3 {
		t.Errorf("invalid len of %s, got %d", x.String(), length)
	}
	x.Add(1)
	x.Add(144)
	x.Add(9)
	length = x.Len()
	if length != 3 {
		t.Errorf("invalid len of %s, got %d", x.String(), length)
	}
	x.Add(98)
	x.Add(99)
	length = x.Len()
	if length != 5 {
		t.Errorf("invalid len of %s, got %d", x.String(), length)
	}
}

func TestCopy(t *testing.T) {
	// TODO: improve test
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	y := x.Copy()
	xLen, yLen := x.Len(), y.Len()
	if xLen != yLen {
		t.Errorf("invalid len of %s and %s", x.String(), y.String())
	}
	for i := range x.words {
		if x.words[i] != y.words[i] {
			t.Errorf("value x: %d but value y: %s", x.words[i], y.words[i])
		}
	}
	x.Add(10)
	if yLen != y.Len() {
		t.Errorf("invalid len %d of y, expected %d", y.Len(), yLen)
	}
}
