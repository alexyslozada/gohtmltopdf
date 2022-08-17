package gohtmltopdf

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

const (
	Executable     = "wkhtmltopdf"
	PlaceHolderArg = "-"
)

type Generator struct {
	stdIn  *bytes.Buffer
	stdOut bytes.Buffer
	stdErr bytes.Buffer
}

func NewGenerator(data *bytes.Buffer) Generator {
	return Generator{
		stdIn:  data,
		stdOut: bytes.Buffer{},
		stdErr: bytes.Buffer{},
	}
}

func (g Generator) run(ctx context.Context) ([]byte, error) {
	// The wkhtmltopdf executable needs to know the source and destination, we can use `-`
	// for stdin and stdout. Then we handle the stdin/stdout to save in memory the process.
	cmd := exec.CommandContext(ctx, Executable, PlaceHolderArg, PlaceHolderArg)
	cmd.Stdin = g.stdIn
	cmd.Stderr = &g.stdErr
	cmd.Stdout = &g.stdOut

	err := cmd.Run()
	if err != nil {
		ctxErr := ctx.Err()
		if ctxErr != nil {
			return nil, ctxErr
		}

		if &g.stdErr != nil {
			errStr := g.stdErr.String()
			if strings.TrimSpace(errStr) != "" {
				return nil, errors.New(errStr)
			}
		}

		return nil, err
	}

	return g.stdOut.Bytes(), nil
}

// writeFile is only for test proposes. --*Don´t use it*--
func writeFile(data []byte) error {
	return ioutil.WriteFile("test.pdf", data, 0666)
}
