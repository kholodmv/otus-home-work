package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		fmt.Println("Flags -from and -to must not be empty")
		return
	}

	fromAbs, err := filepath.Abs(from)
	if err != nil {
		fmt.Println("Error resolving -from path:", err)
		return
	}
	toAbs, err := filepath.Abs(to)
	if err != nil {
		fmt.Println("Error resolving -to path:", err)
		return
	}

	if fromAbs == toAbs {
		fmt.Println("Flags -from and -to must not be equal")
		return
	}

	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
	}
}
