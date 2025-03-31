package xprint

import (
	"unsafe"

	reflect "github.com/goccy/go-reflect"
)

func (p *printer) fmtPointer(value any, verb rune) {
	var u uintptr
	if value == nil && p.verb != 'v' {
		p.buf.writeString(nilAngleString)
		return
	}
	unknownType := false
	switch v := value.(type) {
	case unsafe.Pointer:
		if v == nil {
			p.buf.writeString(nilAngleString)
			return
		}
		u = uintptr(v)
	case uintptr:
		u = uintptr(v)
	case *string:
		u = uintptr(unsafe.Pointer(v))
	case *int:
		u = uintptr(unsafe.Pointer(v))
	case *int8:
		u = uintptr(unsafe.Pointer(v))
	case *int16:
		u = uintptr(unsafe.Pointer(v))
	case *int32:
		u = uintptr(unsafe.Pointer(v))
	case *int64:
		u = uintptr(unsafe.Pointer(v))
	case *uint:
		u = uintptr(unsafe.Pointer(v))
	case *uint8:
		u = uintptr(unsafe.Pointer(v))
	case *uint16:
		u = uintptr(unsafe.Pointer(v))
	case *uint32:
		u = uintptr(unsafe.Pointer(v))
	case *uint64:
		u = uintptr(unsafe.Pointer(v))
	case *uintptr:
		u = uintptr(unsafe.Pointer(v))
	case *float32:
		u = uintptr(unsafe.Pointer(v))
	case *float64:
		u = uintptr(unsafe.Pointer(v))
	case *complex64:
		u = uintptr(unsafe.Pointer(v))
	case *complex128:
		u = uintptr(unsafe.Pointer(v))
	case *bool:
		u = uintptr(unsafe.Pointer(v))
	case *struct{}:
		if verb == 'v' {
			p.buf.writeByte('&')
			p.buf.writeByte('{')
			p.buf.writeByte('}')
			return
		}
		u = uintptr(unsafe.Pointer(v))
	case *interface{}:
		u = uintptr(unsafe.Pointer(v))
	case reflect.Value:

	default:
		unknownType = true
	}

	typ := reflect.TypeOf(value)
	isStructPtr := typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct
	ptrIsNil := u == 0

	switch {
	case unknownType && verb != 's' && verb != 'p' && verb != 'v':
		p.buf.writeString(nilParenString)
		return
	case isStructPtr && verb == 'v':
		p.buf.writeByte('&')
		p.buf.writeByte('{')
		// Print fields here
		p.buf.writeByte('}')
		return

	case isStructPtr && verb == 's':
		p.buf.writeString(typ.String())
		return

	case ptrIsNil && verb == 'v':
		p.buf.writeString(nilAngleString)
		return

	case ptrIsNil && verb == 's':
		p.buf.writeString(percentBangString) // %!
		p.buf.writeByte('s')                 // s
		p.buf.writeByte('(')                 // (
		p.buf.writeString(typ.String())      // *T
		p.buf.writeByte('=')                 // =
		p.buf.writeString(nilAngleString)    // <nil>
		p.buf.writeByte(')')                 // )
		return
	}

	p.buf.writeByte('0')
	p.buf.writeByte('x')

	// Convert uintptr to hex
	buf := make([]byte, 16)
	i := len(buf)
	for u >= 16 {
		i--
		buf[i] = uintToHexDigits[u&0xF]
		u >>= 4
	}
	i--
	buf[i] = uintToHexDigits[u]
	p.buf.write(buf[i:])
}

func (p *printer) printBool() {
	if b, ok := p.arg.(bool); ok {
		if b {
			p.buf.writeString("true")
		} else {
			p.buf.writeString("false")
		}
		return
	}
	p.buf.writeString(percentBangString)
	p.buf.writeByte('t')
	p.buf.writeString(badVerbString)
}

func (p *printer) printReflectType(arg any) {
	p.buf.writeString(reflect.TypeOf(arg).String())
}

func (p *printer) printComplex(v any, verb rune) {
	p.buf.writeByte('(')
	switch v := v.(type) {
	case complex64:
		p.printFloat32(real(v), verb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat32(imag(v), verb)
	case complex128:
		p.printFloat64(real(v), verb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat64(imag(v), verb)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	p.buf.writeByte('i')
	p.buf.writeByte(')')
}
