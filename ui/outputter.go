package ui

import "io"

type Outputter interface {
	Output(io.Writer, string) error
}
