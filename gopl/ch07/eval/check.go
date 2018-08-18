package eval

import (
	"fmt"
	"strings"
)

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

// Check verifies that variable meets requriments.
func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

// Check verifies that literal meets requriments.
func (literal) Check(vars map[Var]bool) error {
	return nil
}

// Check verifies that unary meets requriments.
func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

// Check verifies that binary meets requriments.
func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

// Check verifies that call meets requriments.
func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

// Check verifies that min meets requriments.
func (m min) Check(vars map[Var]bool) error {
	if len(m.args) < 1 {
		return fmt.Errorf("call to min has %d args, want 1 or more",
			len(m.args))
	}
	for _, arg := range m.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}
