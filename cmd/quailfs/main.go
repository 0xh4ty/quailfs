package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("QuailFS")

	switch {
	case len(os.Args) == 1:
		fmt.Println("Starting Client...")

	case os.Args[1] == "--node":
		fmt.Println("Starting Node...")

	default:
		fmt.Println("Unknown argument:", os.Args[1])
		fmt.Println("Usage: quailfs [--node]")
		os.Exit(1)
	}
}
