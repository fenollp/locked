package main

import (
	"fmt"
	"os"

	"github.com/fenollp/locked/locked"
)

func main() {
	t, err := locked.DecodeFile("Lockfile")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf(">>> %#v\n", t)
}
