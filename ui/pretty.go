package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/russross/blackfriday.v2"
)

type WriteFn func(io.Writer, string, ...interface{})

var (
	writeMajorHeader WriteFn = color.New(color.FgHiBlue, color.Bold, color.Underline).FprintfFunc()
	writeMinorHeader WriteFn = color.New(color.FgHiBlue).FprintfFunc()
	writeCode        WriteFn = color.New(color.FgHiMagenta, color.Italic).FprintfFunc()
	writeText        WriteFn = color.New(color.FgHiWhite).FprintfFunc()
	writeLink        WriteFn = color.New(color.FgHiRed, color.Underline).FprintfFunc()
	writeUnknown     WriteFn = color.New(color.FgHiRed).FprintfFunc()
)

type PrettyOutputter struct {
	writer io.Writer

	isContinuingLine bool
}

func (p *PrettyOutputter) Output(w io.Writer, readme string) error {
	p.writer = w
	root := p.parseReadme(readme)
	if err := p.outputNode(root, 0); err != nil {
		return err
	}

	// Output one final newline in case the README ends without one.
	fmt.Fprintln(w)
	return nil
}

func (p *PrettyOutputter) outputNode(n *blackfriday.Node, indent int) error {
	var writeFn WriteFn
	var skipChild bool
	var newline bool

	switch n.Type {
	case blackfriday.Heading:
		writeFn = writeMajorHeader
		newline = true
		skipChild = true
		indent = n.HeadingData.Level - 1
		if n.HeadingData.Level >= 3 {
			writeFn = writeMinorHeader
		}

	case blackfriday.Code:
		writeFn = writeCode
		newline = strings.Contains(string(n.Literal), "\n")

	case blackfriday.Paragraph:
		fallthrough
	case blackfriday.Text:
		writeFn = writeText

	case blackfriday.Link:
		writeFn = writeLink
		skipChild = true

	default:
		writeFn = writeUnknown
	}

	p.write(writeFn, p.nodeContents(n), newline, indent)

	if n.FirstChild != nil && !skipChild {
		if err := p.outputNode(n.FirstChild, indent); err != nil {
			return err
		}
	}
	if n.Next != nil {
		if err := p.outputNode(n.Next, indent); err != nil {
			return err
		}
	}

	return nil
}

func (p *PrettyOutputter) nodeContents(n *blackfriday.Node) string {
	switch n.Type {
	case blackfriday.Heading:
		// The following is required when the root element of a header is a non-text type,
		// for instance a code block:
		// # `code`
		if n.FirstChild.Next != nil {
			return p.nodeContents(n.FirstChild.Next)
		}
		return p.nodeContents(n.FirstChild)

	case blackfriday.Paragraph:
		fallthrough
	case blackfriday.Code:
		fallthrough
	case blackfriday.Text:
		return string(n.Literal)

	case blackfriday.Link:
		text := p.nodeContents(n.FirstChild)
		if len(text) > 0 {
			return fmt.Sprintf("%v <%s>", text, n.LinkData.Destination)
		}
		return fmt.Sprintf("<%s>", n.LinkData.Destination)

	case blackfriday.Image:
		return ""
	case blackfriday.Document:
		return ""

	default:
		return fmt.Sprintf("Type=%v\n", n.Type)
	}
}

func (p *PrettyOutputter) write(w WriteFn, s string, newline bool, indent int) {
	var pad string
	for i := 0; i < indent; i++ {
		pad = pad + "  "
	}

	// Only output a newline if we're not already on one
	if newline && p.isContinuingLine {
		fmt.Fprintln(p.writer, "\n")
		p.isContinuingLine = false
	}

	if !p.isContinuingLine {
		fmt.Fprint(p.writer, pad)
	}

	s = strings.Replace(s, "\n", "\n"+pad, -1)
	s = strings.Replace(s, "%", "%%", -1)
	w(p.writer, s)

	if newline {
		fmt.Fprintln(p.writer, "\n")
	}

	p.isContinuingLine = !newline || strings.HasSuffix(s, "\n")
}

func (*PrettyOutputter) parseReadme(readme string) *blackfriday.Node {
	m := blackfriday.New()
	return m.Parse([]byte(readme))
}
