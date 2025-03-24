package xprint

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

// Printf formats according to a format specifier and returns the resulting string
func Printf(format string, args ...any) string {
	// Fast path for no arguments - just return the format string as-is
	if len(args) == 0 {
		return format
	}

	// Fast path for simple "%s" formatting with string arguments
	if onlyContainsStringPlaceholders(format) && allArgsAreStringLike(args) {
		return fastStringFormat(format, args)
	}

	p := newPrinter()
	p.printf(format, args)
	s := string(p.buf)
	p.free()
	return s
}

func Sprintf(format string, args ...any) string {
	return Printf(format, args...)
}

func Fprintf(w io.Writer, format string, args ...any) (int, error) {
	// Fast path for no arguments - just write the format string as-is
	if len(args) == 0 {
		n, err := w.Write([]byte(format))
		if err != nil {
			return n, errors.New("xprint: Fprint error, provided io.Writer errror: " + err.Error())
		}
		return n, nil
	}

	p := newPrinter()
	p.printf(format, args)
	n, err := w.Write(p.buf)
	p.free()
	if err != nil {
		return n, errors.New("xprint: " + err.Error())
	}
	return n, nil
}

func BigCocncat(s ...string) string {
	var b strings.Builder
	for _, s := range s {
		b.WriteString(s)
	}
	return b.String()
}

var ErrInvalidBool = errors.New("xprint: cant print as bool")

func PrintBool(b bool) string {
	switch b {
	case true:
		return "true"
	default:
		return "false"
	}
}

func PrintInt(i int) string {
	return strconv.Itoa(i)
}
