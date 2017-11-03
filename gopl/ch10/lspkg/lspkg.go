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

func deps(name string) []string {
	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command(binary, "list", "-json", name).Output()
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
	ok := seen[name]
	if ok {
		return
	}
	seen[name] = true
	for _, pkg := range deps(name) {
		walk(seen, pkg)
	}
}

func main() {
	flag.Parse()
	packages := flag.Args()
	_ = packages
	seen := make(map[string]bool)
	walk(seen, packages[0])
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	out = sort.StringSlice(out)
	fmt.Println(strings.Join(out, "\n"))
}
