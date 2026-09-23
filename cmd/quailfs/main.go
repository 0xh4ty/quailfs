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

	case os.Args[1] == "--node" && len(os.Args) > 2 && os.Args[2] == "--relay":
		fmt.Println("Starting Node with Relay...")
		app.RunNode(true)

	case os.Args[1] == "--node":
		fmt.Println("Starting Node...")
		app.RunNode(false)

	case os.Args[1] == "--generate-mnemonic":
		fmt.Println("Generating Mnemonic")
		keys.GenerateMnemonic()

	default:
		fmt.Println("Unknown argument:", os.Args[1])
		fmt.Println("Usage: quailfs [--node] [--relay]")
		os.Exit(1)
	}
}
