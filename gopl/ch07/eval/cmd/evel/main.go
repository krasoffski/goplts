package main

import (
	"fmt"
	"log"

	"github.com/krasoffski/goplts/gopl/ch07/eval"
)

func init() {
	log.SetPrefix("eval: ")
}

func main() {

	expr, err := eval.Parse("min(sqrt(A), pow(y, 3))")
	if err != nil {
		log.Fatal(err)
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		log.Fatal(err)
	}
	fmt.Println(expr.Eval(eval.Env{"A": 100, "y": 15}))
}
