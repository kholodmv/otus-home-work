package main

import (
	"fmt"
	"golang.org/x/example/hello/reverse"
)

func main() {
	s := "Hello, OTUS!"
	r := reverse.String(s)
	fmt.Print(r)
}
