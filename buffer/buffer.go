package buffer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Writer interface {
	io.Writer
	io.StringWriter
}

type Reader interface {
	io.Reader
}

type Buffer struct {
	Output Writer
	Input  Reader

	name       string
	beforeLine string
	close      bool
}

func NewBuffer(name string) Buffer {
	return Buffer{
		Input:      os.Stdin,
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
	n := 0
	NewLine := []byte("\n")
	splited := bytes.Split(p, NewLine)
	for _, split := range splited {
		bytes.Trim(split, "\r")
	}

	lastIndex := len(splited) - 1
	for i := 0; i < lastIndex; i++ {
		n += b.writeLine(splited[i]) + len(NewLine)
		b.Output.Write(NewLine)
	}
	if len(splited[lastIndex]) != 0 {
		n += b.writeLine(splited[lastIndex])
	}

	return n, nil
}

func (b *Buffer) writeLine(line []byte) int {
	line_len := len(line)
	if strings.Contains(b.beforeLine, "%v") {
		prefix := fmt.Sprintf(b.beforeLine, b.name)
		line_len += len(prefix)
		b.Output.WriteString(prefix)
	} else if len(b.beforeLine) > 0 {
		line_len += len(b.beforeLine)
		b.Output.WriteString(b.beforeLine)
	}
	b.Output.Write(line)

	return line_len
}

func (b *Buffer) WriteString(s string) (int, error) {
	return b.Write([]byte(s))
}

func (b *Buffer) Read(p []byte) (int, error) {
	return b.Input.Read(p)
}

func (b *Buffer) Close(p []byte) error {
	if b.close {
		return errors.New("you can't reclose")
	}
	if o, ok := b.Output.(io.Closer); ok {
		o.Close()
	}
	if i, ok := b.Input.(io.Closer); ok {
		i.Close()
	}
	b.close = true
	return nil
}
