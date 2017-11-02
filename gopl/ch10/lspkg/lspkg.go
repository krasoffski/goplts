package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func deps(name string) []string {
	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command(binary, "list", "-json").Output()
	if err != nil {
		log.Fatal(err)
	}
	properties := struct{ Deps []string }{}
	if err := json.Unmarshal(out, &properties); err != nil {
		log.Fatal(err)
	}
	return properties.Deps
}

func walk(seen map[string]bool, name string) {

}

func main() {
	flag.Parse()
	packages := flag.Args()
	_ = packages
	// var seen map[string]bool

	fmt.Println(deps(packages[0]))
}
