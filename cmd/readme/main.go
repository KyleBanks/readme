package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Invalid usage, expected 'readme username/repository'")
	}

	repo := os.Args[1]
	fmt.Printf("github.com/%v/README.md", repo)
}
