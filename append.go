package xprint

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
	switch {
	case len(items) == 0:
		return b
	case len(items) == 1 && items[0] == nil:
		b = append(b, nilAngleString...)
		return b
	}
	p := newPrinter()
	p.printf(format, items)

	b = append(b, p.buf...)
	p.free()
	return b
}

// doappend formats the arguments using their default formats (%v verb)
// and places them into p.buf with appropriate spacing.
func (p *printer) doappend(items []any) {
	p.print(items)
}
