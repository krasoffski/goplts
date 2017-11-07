package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

var binaryPath string

func init() {
	var err error
	binaryPath, err = exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
}

type properties struct {
	Name       string
	ImportPath string
	Deps       []string
}

func getProperties(name string) (*properties, error) {
	out, err := exec.Command(binaryPath, "list", "-json", name).Output()
	if err != nil {
		return nil, err
	}
	var p properties
	if err := json.Unmarshal(out, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func main() {
	flag.Parse()
	packages := flag.Args()
	if len(packages) <= 0 {
		log.Fatal("please, specify package name")
	}
	dependencies := make(map[string]int)
	for _, p := range packages {
		prop, err := getProperties(p)
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range prop.Deps {
			dependencies[d]++
		}
	}
	fmt.Println(dependencies)
}
