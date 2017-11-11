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
	if binaryPath, err = exec.LookPath("go"); err != nil {
		log.Fatal(err)
	}
}

func getPackages(workspace string) ([]string, error) {
	// dir := filepath.Join(workspace, "...") // does not work with '...'
	dir := filepath.Clean(workspace) + string(filepath.Separator) + "..."
	out, err := exec.Command(binaryPath, "list", dir).Output()
	if err != nil {
		return nil, err
	}
	// TODO: improve this.
	str := strings.TrimSpace(string(out))
	return strings.Split(str, "\n"), nil
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

func getUsedPackages(deps map[string]bool, packages []string) []string {
	used := make([]string, 0)
	for _, pkg := range packages {
		if !deps[pkg] {
			continue
		}
		used = append(used, pkg)
	}
	return used
}

func unique(in []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(in)) // think about in-place modification
	for _, item := range in {
		if seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func main() {
	wsPath := flag.String("workspace", ".", "workspace relative path")
	flag.Parse()
	args := flag.Args()
	if len(args) <= 0 {
		log.Fatal("please, specify package names")
	}

	argPkgs := unique(args)

	wsPkgs, err := getPackages(*wsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range wsPkgs {
		deps, err := getDependencies(p)
		if err != nil {
			log.Fatal(err)
		}
		used := getUsedPackages(deps, argPkgs)
		if len(used) > 0 {
			fmt.Printf("%s: [%s]\n", p, strings.Join(used, ", "))
		}
	}
}
