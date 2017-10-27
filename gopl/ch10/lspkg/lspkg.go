package main

import (
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
	fmt.Println(string(out))
}
