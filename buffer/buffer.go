package buffer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

type Writer interface {
	io.Writer
	io.StringWriter
}

type Buffer struct {
	Output Writer

	name       string
	beforeLine string
	close      bool
}

func NewBuffer(name string) Buffer {
	return Buffer{
		Output:     os.Stdout,
		name:       name,
		beforeLine: "[%v]: ",
		close:      false,
	}
}

func (b *Buffer) Write(p []byte) (int, error) {
	if b.close {
		return 0, errors.New("you can't write want the buffer is closed")
	}
	splited := bytes.Split(p, NewLine)
	for i := 0; i < len(splited)-1; i++ {
		line := splited[i]
		b.Output.WriteString(fmt.Sprintf(b.beforeLine, b.name))
		b.Output.Write(line)
		b.Output.Write(NewLine)
	}
	return len(p), nil
}

func (b *Buffer) Close(p []byte) error {
	if b.close {
		return errors.New("you can't reclose")
	}
	if o, ok := b.Output.(io.Closer); ok {
		o.Close()
	}
	b.close = true
	return nil
}
