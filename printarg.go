package xprint

import (
	"strconv"

	reflect "github.com/goccy/go-reflect"
)

// printArg formats arg in the manner specified by the verb
// and appends it to p.buf.
func (p *printer) printArg() {
	// Handle nil
	if p.arg == nil {
		switch p.verb {
		case 'T', 'v':

			p.buf.writeString(nilString)
		default:

			p.buf.writeNilArg(p.verb)
		}
		return
	}
	// Handle based on type and verb
	switch p.verb {
	case 'T':
		p.printReflectType(p.arg)
		return
	case 't':
		p.printBool(p.arg)
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(p.arg), p.verb)
	}

	// Handle by type
	switch v := p.arg.(type) {
	case []byte:
		p.buf = append(p.buf, v...)
	case string:
		if p.fmt.widPresent && p.verb == 's' {
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
		switch p.verb {
		case 't':
			p.printBool(v)
		case 's':
			boolstr := percentBangString + "s(" + "bool" + "=" + strconv.FormatBool(v) + ")"
			p.buf = append(p.buf, boolstr...)
		}
	case int:
		p.fmtInt()

	case int8:
		p.fmtInt8()

	case int16:
		p.fmtInt16()

	case int32:
		p.fmtInt32()

	case int64:
		p.fmtInt64()

	case uint:
		p.fmtUint()

	case uint8:
		p.fmtUint8()

	case uint16:
		p.fmtUint16()

	case uint32:
		p.fmtUint32()

	case uint64:
		p.fmtUint64()

	case uintptr:
		p.fmtUintptr()

	case float32:
		// If precision is explicitly specified, use printFloat
		// Otherwise use our specialized formatter with proper defaults
		if p.fmt.precPresent {
			p.printFloat(v, p.verb)
		} else {
			p.printFloat32(v, p.verb)
		}
	case float64:
		// If precision is explicitly specified, use printFloat
		// Otherwise use our specialized formatter with proper defaults
		if p.fmt.precPresent {
			p.printFloat(v, p.verb)
		} else {
			p.printFloat64(v, p.verb)
		}
	case complex64, complex128:
		p.printComplex(v, p.verb)
	default:
		if p.handleMethods(p.verb) {
			return
		}

		p.value = reflect.ValueOf(p.arg)
		p.printValue(p.value, p.verb, 0)
	}
}

func (p *printer) handleMethods(verb rune) bool {
	if err, ok := p.arg.(error); ok {
		func() {
			defer func() {
				if r := recover(); r != nil {
					p.catchPanic(r, verb, "Error")
				}
			}()
			p.fmt.fmtString(err.Error())
		}()
		return true
	}
	if stringer, ok := p.arg.(Stringer); ok {
		func() {
			defer func() {
				if r := recover(); r != nil {
					p.catchPanic(r, verb, "Stringer")
				}
			}()
			p.fmt.fmtString(stringer.String())
		}()
		return true
	}
	return false
}
