package xprint

import "strconv"

// printFloat64 handles float64 formatting with default precision
func (p *printer) printFloat64(v float64, verb rune) {
	switch verb {
	case 'f', 'e', 'E':
		// Default precision 6
		p.buf.writeString(strconv.FormatFloat(v, byte(verb), 6, 64))
	case 'F':
		// Special case: F uses 'f' format with precision 6
		p.buf.writeString(strconv.FormatFloat(v, 'f', 6, 64))
	case 'v':
		// Special case: %v uses 'g' format
		p.buf.writeString(strconv.FormatFloat(v, 'g', -1, 64))
	default:
		// 'g', 'G', 'b', 'x', 'X' - default precision -1
		p.buf.writeString(strconv.FormatFloat(v, byte(verb), -1, 64))
	}
}

// printFloat32 handles float32 formatting with default precision
func (p *printer) printFloat32(v float32, verb rune) {
	switch verb {
	case 'f', 'e', 'E':
		// Default precision 6
		p.buf.writeString(strconv.FormatFloat(float64(v), byte(verb), 6, 32))
	case 'F':
		// Special case: F uses 'f' format with precision 6
		p.buf.writeString(strconv.FormatFloat(float64(v), 'f', 6, 32))
	case 'v':
		// Special case: %v uses 'g' format
		p.buf.writeString(strconv.FormatFloat(float64(v), 'g', -1, 32))
	default:
		// 'g', 'G', 'b', 'x', 'X' - default precision -1
		p.buf.writeString(strconv.FormatFloat(float64(v), byte(verb), -1, 32))
	}
}
