package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
)

func imports(name string) []string {
	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command(binary, "list", "-json", name).Output()
	if err != nil {
		log.Fatal(err)
	}
	properties := struct{ Imports []string }{}
	if err := json.Unmarshal(out, &properties); err != nil {
		log.Fatal(err)
	}
	return properties.Imports
}

func walk(seen map[string]bool, name string) {
	ok := seen[name]
	if ok {
		return
	}
	seen[name] = true
	for _, pkg := range imports(name) {
		walk(seen, pkg)
	}
}

func main() {
	flag.Parse()
	packages := flag.Args()
	_ = packages
	deps := make(map[string]bool)
	walk(deps, packages[0])
	out := make([]string, 0, len(deps))
	for k := range deps {
		out = append(out, k)
	}
	out = sort.StringSlice(out)
	fmt.Println(strings.Join(out, "\n"))
}
