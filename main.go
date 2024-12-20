package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smol-go/smol-git/internal/repository"
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
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: smolgit add <file>")
			os.Exit(1)
		}
		if err := handleAdd(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "status":
		if err := handleStatus(); err != nil {
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

	if _, err := repository.Init(path); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	fmt.Println("Initialized empty Git repository")
	return nil
}

func handleAdd(path string) error {
	repo, err := repository.Open(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	if err := repo.Add(path); err != nil {
		return fmt.Errorf("failed to add file: %w", err)
	}

	fmt.Printf("Added %s to staging area\n", path)
	return nil
}

func handleStatus() error {
	repo, err := repository.Open(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	status, err := repo.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fmt.Println(status)
	return nil
}
