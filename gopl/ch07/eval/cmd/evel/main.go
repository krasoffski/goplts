package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/krasoffski/goplts/gopl/ch07/eval"
)

func init() {
	log.SetPrefix("eval: ")
	log.SetFlags(0)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter formula: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	// expr, err := eval.Parse("min(sqrt(A), pow(y, 3))")
	expr, err := eval.Parse(text)
	if err != nil {
		log.Fatalln(err)
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		log.Fatal(err)
	}
	env := make(eval.Env)
	for v := range vars {
		var val float64
		fmt.Printf("Enter %s: ", v)
		_, err := fmt.Scanf("%f\n", &val)
		if err != nil {
			log.Fatalln(err)
		}
		env[v] = val
	}
	fmt.Println(expr.Eval(env))
}
