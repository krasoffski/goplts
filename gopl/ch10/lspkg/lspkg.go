package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

func getPackages(workspace string) ([]string, error) {
	dir := filepath.Join(workspace, "...")
	out, err := exec.Command(binaryPath, "list", dir).Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}

func getDependencies(name string) (map[string]bool, error) {
	out, err := exec.Command(binaryPath, "list", "-json", name).Output()
	if err != nil {
		return nil, err
	}
	var p struct{ Deps []string }
	if err := json.Unmarshal(out, &p); err != nil {
		return nil, err
	}
	deps := make(map[string]bool)
	for _, pkg := range p.Deps {
		deps[pkg] = true
	}
	return deps, nil
}

func main() {
	workspace := flag.String("ws", "", "workspace path")
	flag.Parse()
	lsPackages := flag.Args()
	if len(lsPackages) <= 0 {
		log.Fatal("please, specify package names")
	}
	wsPackages, err := getPackages(*workspace)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range wsPackages {
		deps, err := getDependencies(p)
		if err != nil {
			log.Fatal(err)
		}
		for _, lsPkg := range lsPackages {
			if deps[lsPkg] {
				fmt.Println(p)
				break
			}
		}
	}
}
