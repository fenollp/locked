package main

import (
	"fmt"

	"github.com/fenollp/locked/locked"
)

func main() {
	t, err := locked.DecodeFile("Lockfile")
	if err != nil {
		panic(err)
	}
	fmt.Printf(">>> %#v\n", t)
}
