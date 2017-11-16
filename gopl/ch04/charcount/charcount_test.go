package charcount

import (
	"strings"
	"testing"
)

func fromString(str string) (map[string]int, error) {
	return Counter(strings.NewReader(str))
}

func TestEmptyInputErr(t *testing.T) {
	if _, err := fromString(""); err != nil {
		t.Errorf("expected nil but got '%s'", err)
	}
}

func TestEmptyInputResult(t *testing.T) {
	if res, _ := fromString(""); len(res) != 0 && res != nil {
		t.Errorf("expected empty map but got '%v'", res)
	}
}
