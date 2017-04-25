package eval

import (
	"bytes"
	"fmt"
)

// An Expr is an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]bool) error
	// String returns string representation of expression.
	String() string
}

// A Var identifies a variable, e.g., x.
type Var string

func (e Var) String() string {
	return string(e)
}

// A literal is numeric constant.
type literal float64

func (e literal) String() string {
	return fmt.Sprintf("%g", e)
}

type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

func (e unary) String() string {
	return fmt.Sprintf("(%c%s)", e.op, e.x)
}

// A binary represents a binary operator expression.
type binary struct {
	op   rune
	x, y Expr
}

func (e binary) String() string {
	return fmt.Sprintf("(%s%c%s)", e.x, e.op, e.y)
}

// A call represents a function call expression, e.g., six(x).
type call struct {
	fn   string
	args []Expr
}

func (e call) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s(", e.fn)
	for i, arg := range e.args {
		if i > 0 {
			buf.WriteString(", ")
		}
		write(buf, arg)
	}
	buf.WriteByte(')')
	return buf.String()
}

// A min computes the minimum from it's arguments.
type min struct {
	args []Expr
}

func (m min) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "min(")
	for i, arg := range m.args {
		if i > 0 {
			buf.WriteString(", ")
		}
		write(buf, arg)
	}
	buf.WriteByte(')')
	return buf.String()
}
