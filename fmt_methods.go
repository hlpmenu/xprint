package xprint

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// Core formatting methods
func (f *fmt) fmtBool(v bool) {
	if v {
		f.buf.writeString("true")
	} else {
		f.buf.writeString("false")
	}
}

func (f *fmt) fmtString(s string) {
	f.buf.writeString(s)
}

func (f *fmt) fmtBytes(v []byte) {
	f.buf.write(v)
}

func (f *fmt) fmtFloat(v float64, size int, verb rune, prec int) {
	// Handle sign
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	} else if f.plus {
		sign = "+"
	} else if f.space {
		sign = " "
	}

	// Set default precision if not specified
	if prec < 0 {
		switch verb {
		case 'e', 'E', 'f', 'F':
			prec = 6
		case 'g', 'G':
			prec = -1 // Use shortest representation
		default:
			prec = 6
		}
	}

	// Format based on verb
	var format byte
	switch verb {
	case 'f', 'F':
		format = 'f'
	case 'e', 'E':
		format = byte(verb)
	case 'g', 'G':
		format = byte(verb)
	default:
		format = 'g'
	}

	// Convert to string
	num := strconv.FormatFloat(v, format, prec, size)

	// Handle special cases for 'g'/'G' format
	if (verb == 'g' || verb == 'G') && f.sharp && strings.IndexByte(num, '.') < 0 {
		num += "."
	}

	// Combine all parts
	s := sign + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) fmtInt(v int64, base int, verb rune) {
	// Build the string
	var s string

	// Determine sign using switch
	sign := ""
	switch {
	case v < 0:
		sign = "-"
		v = -v
	case f.plus:
		sign = "+"
	case f.space:
		sign = " "
	}

	// Convert number to string
	num := strconv.FormatInt(v, base)
	if base == 16 && verb == 'X' {
		num = strings.ToUpper(num)
	}

	// Add prefix if requested
	prefix := ""
	if f.sharp {
		switch base {
		case 2:
			prefix = "0b"
		case 8:
			prefix = "0o"
		case 16:
			prefix = "0x"
		}
	}

	// Combine all parts
	s = sign + prefix + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) fmtUint(v uint64, base int, verb rune) {
	// Build the string
	var s string
	sign := ""
	if f.plus {
		sign = "+"
	} else if f.space {
		sign = " "
	}

	// Convert number to string
	num := strconv.FormatUint(v, base)
	if base == 16 && verb == 'X' {
		num = strings.ToUpper(num)
	}

	// Add prefix if requested
	prefix := ""
	if f.sharp {
		switch base {
		case 2:
			prefix = "0b"
		case 8:
			prefix = "0o"
		case 16:
			prefix = "0x"
		}
	}

	// Combine all parts
	s = sign + prefix + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) pad(s string) {
	if !f.widPresent || f.wid == 0 {
		f.buf.writeString(s)
		return
	}

	width := f.wid - utf8.RuneCountInString(s)
	if width <= 0 {
		f.buf.writeString(s)
		return
	}

	// Left padding
	if !f.minus {
		if f.zero && !strings.ContainsAny(s, ".-+") {
			//nolint:all
			for i := 0; i < width; i++ {
				f.buf.writeByte('0')
			}
		} else {
			//nolint:all
			for i := 0; i < width; i++ {
				f.buf.writeByte(' ')
			}
		}
	}

	f.buf.writeString(s)

	// Right padding
	if f.minus {
		//nolint:all
		for i := 0; i < width; i++ {
			f.buf.writeByte(' ')
		}
	}
}
