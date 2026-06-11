package iostreams

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

type IOStreams struct {
	Out     io.Writer
	ErrOut  io.Writer
	In      io.Reader
	IsTTY   bool
	Output  string
	JQ      string
	Columns []string
}

func System() *IOStreams {
	return &IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
		In:     os.Stdin,
		IsTTY:  isTerminal(),
	}
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
