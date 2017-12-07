package ui

import (
	"fmt"
	"io"
)

type RawOutputter struct{}

func (RawOutputter) Output(w io.Writer, readme string) error {
	fmt.Fprintln(w, readme)
	return nil
}
