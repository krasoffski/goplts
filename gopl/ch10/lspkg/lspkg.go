package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func deps(name string) []string {

}

func walk(seen map[string]bool, name string) {

}

func main() {
	flag.Parse()
	packages := flag.Args()
	_ = packages
	var seen map[string]bool

	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command(binary, "list", "-json").Output()
	if err != nil {
		log.Fatal(err)
	}
	deps := struct{ Deps []string }{}
	if err := json.Unmarshal(out, &deps); err != nil {
		log.Fatal(err)
	}
	fmt.Println(deps.Deps)
}
