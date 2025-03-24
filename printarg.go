package xprint

import (
	"reflect"
	"strconv"
)

// printArg formats arg in the manner specified by the verb
// and appends it to p.buf.
func (p *printer) printArg(arg any, verb rune) {
	// Handle nil
	if arg == nil {
		p.buf.writeString(nilAngleString)
		return
	}

	// Handle based on type and verb
	switch verb {
	case 'T':
		p.printReflectType(arg)
		return
	case 't':
		p.printBool(arg)
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), verb)
	}
	// Handle by type
	switch v := arg.(type) {
	case []byte:
		p.buf.writeByte('[')
		for i, b := range v {
			if i > 0 {
				p.buf.writeByte(' ')
			}
			p.printInt(b, 10, 'd')
		}
		p.buf.writeByte(']')
	case string:
		p.buf = append(p.buf, v...)
	case bool:
		switch verb {
		case 't', 'v':
			p.printBool(v)
		case 's':
			boolstr := percentBangString + "s(" + "bool" + "=" + strconv.FormatBool(v) + ")"
			p.buf = append(p.buf, boolstr...)
		}
	case int, int8, int16, int32, int64:
		p.printInt(v, 10, verb)
	case uint, uint8, uint16, uint32, uint64:
		p.printInt(v, 10, verb)
	case uintptr:
		// Special handling for uintptr
		if verb == 'v' {
			// For %v, print in decimal
			p.printInt(v, 10, 'd')
		} else if verb == 'x' || verb == 'X' {
			// For %x or %X, print in hex
			p.printInt(v, 16, verb)
		} else {
			// For other verbs, use printInt with decimal base
			p.printInt(v, 10, verb)
		}
	case float32:
		// If precision is explicitly specified, use printFloat
		// Otherwise use our specialized formatter with proper defaults
		if p.fmt.precPresent {
			p.printFloat(v, verb)
		} else {
			p.printFloat32(v, verb)
		}
	case float64:

		// If precision is explicitly specified, use printFloat
		// Otherwise use our specialized formatter with proper defaults
		if p.fmt.precPresent {
			p.printFloat(v, verb)
		} else {
			p.printFloat64(v, verb)
		}
	case complex64, complex128:
		p.printComplex(v, verb)
	default:
		// Check for interface methods like error.Error() or Stringer.String()
		if p.handleMethods(verb) {
			return
		}

		// Store the argument in p.arg for reflection
		p.arg = arg

		// Use reflection for other types
		p.value = reflect.ValueOf(arg)
		p.printValue(p.value, verb, 0)
	}
}
