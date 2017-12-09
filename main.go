package main

import (
	"fmt"
	"io"
	"os"

	"github.com/KyleBanks/readme/git"
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

	resolver, err := args.Resolver()
	check(err)

	readme, err := fetchReadme(resolver)
	check(err)

	out, err := args.Outputter()
	check(err)

	err = output(os.Stdout, readme, out)
	check(err)
}

func fetchReadme(r git.Resolver) (string, error) {
	var err error
	var contents string
	for _, filename := range readmeVariations {
		contents, err = r.Resolve(filename)
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
	fmt.Printf("\t%v: skips rendering of images\n", FlagNoImages)
}

func check(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERROR: %v\n", err)
	os.Exit(1)
}
