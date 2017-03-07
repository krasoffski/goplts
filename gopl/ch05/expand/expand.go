package main

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\$\S+`) // TODO: need find more solid regexp

func expand(s string, f func(string) string) string {
	// NOTE: in some cases (mem) it is better to request new string each time.
	for _, w := range re.FindAllString(s, -1) {
		s = strings.Replace(s, w, f(w[1:]), 1)
	}
	return s
}

func main() {
	str := "$var1 $var2 expand $var3"
	fmt.Println(str)
	str = expand(str, func(s string) string {
		return "${" + s + "}"
	})
	fmt.Println(str)
	str = expand(str, strings.ToUpper)
	fmt.Println(str)
}
