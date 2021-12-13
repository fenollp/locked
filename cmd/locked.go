package main

import (
	"fmt"
	"os"

	"github.com/fenollp/locked/locked"
)

func main() {
	if err := locked.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
