package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smol-go/smol-git/internal/respository"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: smolgit <command> [arguments]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "init":
		if err := handleInit(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func handleInit() error {
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if _, err := respository.Init(path); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	fmt.Println("Initialized empty Git repository")
	return nil
}
