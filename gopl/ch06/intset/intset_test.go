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
