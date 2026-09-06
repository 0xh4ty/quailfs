package main

import (
	"fmt"
	"github.com/0xh4ty/quailfs/internal/app"
	"os"
)

func main() {
	fmt.Println("QuailFS")

	switch {
	case len(os.Args) == 1:
		fmt.Println("Starting Client...")
		app.RunClient()

	case os.Args[1] == "--node":
		fmt.Println("Starting Node...")
		app.RunNode()

	default:
		fmt.Println("Unknown argument:", os.Args[1])
		fmt.Println("Usage: quailfs [--node]")
		os.Exit(1)
	}
}
