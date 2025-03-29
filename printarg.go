package xprint

import (
	"reflect"
	"strconv"
)

// printArg formats arg in the manner specified by the verb
// and appends it to p.buf.
func (p *printer) printArg(arg any, verb rune) {
	// Handle nil
	if arg == nil || p.arg == nil {
		switch verb {
		case 'T', 'v':
			p.buf.writeString(nilString)
		default:
			p.buf.writeNilArg(verb)

		}
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
		p.buf = append(*p.fmt.buf, v...)
	case string:
		if p.fmt.widPresent && verb == 's' {
			width := p.fmt.wid - len(v)
			if width > 0 {
				// Left padding (right-aligned)
				if !p.fmt.minus {
					for i := 0; i < width; i++ {
						p.buf.writeByte(' ')
					}
				}
				// Write the string
				p.buf.writeString(v)
				// Right padding (left-aligned)
				if p.fmt.minus {
					for i := 0; i < width; i++ {
						p.buf.writeByte(' ')
					}
				}
			} else {
				// Width is less than string length, just write string
				p.buf.writeString(v)
			}
		} else {
			p.buf = append(p.buf, v...)
		}
	case bool:
		switch verb {
		case 't':
			p.printBool(v)
		case 's':
			boolstr := percentBangString + "s(" + "bool" + "=" + strconv.FormatBool(v) + ")"
			p.buf = append(p.buf, boolstr...)
		}
	case int, int8, int16, int32, int64:
		p.printInt(v, verb)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		p.printInt(v, verb)
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
		if p.handleMethods(verb) {

			return
		}

		p.value = reflect.ValueOf(p.arg)

		p.printValue(p.value, verb, 0)
	}

}
