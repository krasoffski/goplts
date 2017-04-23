package eval

// An Expr is an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]bool) error
}

// A Var identifies a variable, e.g., x.
type Var string

// A literal is numeric constant.
type literal float64

type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

// A binary represents a binary operator expression.
type binary struct {
	op   rune
	x, y Expr
}

// A call represents a function call expression, e.g., six(x).
type call struct {
	fn   string
	args []Expr
}
