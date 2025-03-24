package xprint

import (
	fmtpkg "fmt"
	"reflect"
)

// Append formats using the default formats for its operands, appends the result to
// the byte slice, and returns the updated slice.
func Append(b []byte, items ...any) []byte {
	// Fast path for no arguments - just return the input as-is
	p := newPrinter()
	p.doappend(items)
	b = append(b, p.buf...)
	p.free()
	return b
}

// Appendf formats according to a format specifier, appends the result to
// the byte slice, and returns the updated slice.
func Appendf(b []byte, format string, items ...any) []byte {
	// Fast path for no arguments - just return the input as-is

	p := newPrinter()
	p.printf(format, items)
	b = append(b, p.buf...)
	p.free()
	return b
}

// doappend formats the arguments using their default formats (%v verb)
// and places them into p.buf with appropriate spacing.
func (p *printer) doappend(items []any) {
	prevString := false
	for argNum, arg := range items {
		isString := arg != nil && reflect.TypeOf(arg).Kind() == reflect.String
		// Add a space between two non-string arguments
		if argNum > 0 && !isString && !prevString {
			p.buf.writeByte(' ')
		}
		p.printArg(arg, 'v')
		prevString = isString
	}
}

func ref() {
	fmtpkg.Appendf(nil, "")
}
