package buffer

import (
	"fmt"
	"os"
)

var Stdout Buffer = Buffer{
	name:   "stdout",
	Output: os.Stdout,
}

var Stderr Buffer = Buffer{
	name:   "stderr",
	Output: os.Stderr,
}

var Stdin Buffer = Buffer{
	name:  "stdin",
	Input: os.Stdin,
}

func Print(a ...any) (int, error) {
	return fmt.Fprint(&Stdout, a...)
}
func Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(&Stdout, format, a...)
}
func Println(a ...any) (int, error) {
	return fmt.Fprintln(&Stdout, a...)
}

func Eprint(a ...any) (int, error) {
	return fmt.Fprint(&Stderr, a...)
}
func Eprintf(format string, a ...any) (int, error) {
	return fmt.Fprintf(&Stderr, format, a...)
}
func Eprintln(a ...any) (int, error) {
	return fmt.Fprintln(&Stderr, a...)
}
