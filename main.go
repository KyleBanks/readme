package main

import (
	"fmt"
	"io"
	"os"

	"github.com/KyleBanks/readme/git"
	"github.com/KyleBanks/readme/git/http"
	"github.com/KyleBanks/readme/ui"
)

var (
	readmeVariations = []string{"README.md", "Readme.md", "README.txt", "Readme.txt", "README", "Readme"}
)

func main() {
	args, err := parseArgs(os.Args)
	if err != nil {
		printUsageError(err)
		return
	}

	readme, err := fetchReadme(http.NewGitHubHttpResolver(), args.Username, args.Repository)
	if err != nil {
		fmt.Printf("ERROR: failed to fetch README: %v\n", err)
		os.Exit(1)
		return
	}

	if err := output(os.Stdout, readme, args.Outputter()); err != nil {
		fmt.Printf("ERROR: failed to generate output: %v\n", err)
		os.Exit(1)
		return
	}
}

func fetchReadme(r git.Resolver, username, repository string) (string, error) {
	var err error
	var contents string
	for _, filename := range readmeVariations {
		contents, err = r.Resolve(username, repository, filename)
		if err == nil {
			break
		}
	}

	return contents, err
}

func output(w io.Writer, readme string, out ui.Outputter) error {
	return out.Output(w, readme)
}

func printUsageError(err error) {
	fmt.Printf("ERROR: %v\n\n", err)
	printUsage()
	os.Exit(2)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\treadme username/repository [options...]")
	fmt.Println("\nOptions:")
	fmt.Printf("\t%v: outputs the readme as plain text\n", FlagRaw)
}
