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

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter formula: ")
	text, err := reader.ReadString('\n')
	checkErr(err)

	expr, err := eval.Parse(text)
	checkErr(err)

	vars := make(map[eval.Var]bool)
	err = expr.Check(vars)
	checkErr(err)

	var val float64
	env := make(eval.Env)
	for v := range vars {
		fmt.Printf("Enter %s: ", v)
		_, err = fmt.Scanf("%f\n", &val)
		checkErr(err)
		env[v] = val
	}
	fmt.Println(expr.Eval(env))
}
