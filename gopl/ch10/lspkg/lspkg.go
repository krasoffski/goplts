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

func main() {
	flag.Parse()
	lsPackages := flag.Args()
	if len(lsPackages) <= 0 {
		log.Fatal("please, specify package names")
	}

	wsPackages, err := getPackages(".")
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
				fmt.Printf("%s - %s\n", p, lsPkg)
				break
			}
		}
	}
}
