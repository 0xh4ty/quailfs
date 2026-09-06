package main

import (
	"fmt"
	"github.com/0xh4ty/quailfs/internal/app"
	"github.com/0xh4ty/quailfs/internal/keys"
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

	case os.Args[1] == "--generate-mnemonic":
		fmt.Println("Generating Mnemonic")
		keys.GenerateMnemonic()

	default:
		fmt.Println("Unknown argument:", os.Args[1])
		fmt.Println("Usage: quailfs [--node]")
		os.Exit(1)
	}
}
