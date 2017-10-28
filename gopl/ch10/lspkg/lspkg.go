package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func main() {
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
