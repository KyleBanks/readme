package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KyleBanks/readme/git"
	"github.com/KyleBanks/readme/git/http"
	"github.com/KyleBanks/readme/ui"
)

const (
	FlagRaw = "-raw"
)

var (
	ErrRepositoryNotProvided   = errors.New("Missing repository argument")
	ErrInvalidRepositoryFormat = errors.New("Invalid repository format")
)

type Args struct {
	Username   string
	Repository string

	Raw bool
}

func (a Args) Outputter() ui.Outputter {
	if a.Raw {
		return ui.RawOutputter{}
	}
	return &ui.PrettyOutputter{}
}

func (a Args) Resolver() git.Resolver {
	return http.NewGitHubHttpResolver()
}

func parseArgs(args []string) (*Args, error) {
	if len(args) < 2 {
		return nil, ErrRepositoryNotProvided
	}

	repoTokens := strings.Split(args[1], "/")
	if len(repoTokens) != 2 {
		return nil, ErrInvalidRepositoryFormat
	}

	parsed := Args{
		Username:   repoTokens[0],
		Repository: repoTokens[1],
	}

	if len(args) > 2 {
		for _, arg := range args[2:] {
			switch arg {
			case FlagRaw:
				parsed.Raw = true
			default:
				return nil, fmt.Errorf("Unknown argument: %v", arg)
			}
		}
	}

	return &parsed, nil
}
