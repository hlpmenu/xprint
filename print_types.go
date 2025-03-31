package xprint

import (
	"log"
	"reflect"
	"unsafe"
)

func (p *printer) fmtPointer(value any, verb rune) {
	var u uintptr
	if value == nil && p.verb != 'v' {
		p.buf.writeString(nilAngleString)
		return
	}

	switch v := value.(type) {

	case unsafe.Pointer:
		log.Printf("triggeer unsafe.Pointer")
		u = uintptr(v)
	case uintptr:
		log.Printf("triggeer uintptr")

		u = v
	case *string:

		log.Printf("triggeer *string")
		log.Printf("in *string: %p", reflect.TypeOf(v).String())

		u = uintptr(unsafe.Pointer(v))
	case *int:
		log.Printf("triggeer *int")
		u = uintptr(unsafe.Pointer(v))
	case *int8:
		log.Printf("triggeer *int8")
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
		log.Printf("triggeer *interface{}")
		u = uintptr(unsafe.Pointer(v))
	default:
		log.Printf("default")
		switch verb {
		case 's', 'p', 'v':

		default:
			log.Printf("default nested")
			p.buf.writeString(nilParenString)
			return
		}
	}

	if reflect.TypeOf(value).Kind() == reflect.Ptr && reflect.TypeOf(value).Elem().Kind() == reflect.Struct {
		if verb == 'v' {
			p.buf.writeByte('&')
			p.buf.writeByte('{')
			// Print fields here
			p.buf.writeByte('}')
			return
		} else if verb == 's' {
			p.buf.writeString(reflect.TypeOf(value).String())
			return
		}
	}
	log.Printf("verb before u == 0: %v", verb)
	if u == 0 {

		refv := reflect.ValueOf(value)
		kind := refv.Kind()
		log.Printf("in if kind: %s", kind.String())

		log.Printf("trigger if u == 0 && p.verb == 'u': %v p.verb: %v", u, p.verb)
		switch p.verb {
		case 'v':
			p.buf.writeString(nilAngleString)
			return
		case 's':
			p.buf.writeString(percentBangString)              // %!
			p.buf.writeByte('s')                              // s
			p.buf.writeByte('(')                              // (
			p.buf.writeString(reflect.TypeOf(value).String()) // *string
			p.buf.writeByte('=')                              // =
			p.buf.writeString(nilAngleString)                 // <nil>
			p.buf.writeByte(')')                              // )

			return
		}

	}

	p.buf.writeByte('0')
	p.buf.writeByte('x')
	log.Printf("value: %d", u)
	log.Printf("value: %d", u)

	// Convert uintptr to hex
	const digits = "0123456789abcdef"
	buf := make([]byte, 16)
	i := len(buf)
	for u >= 16 {
		i--
		buf[i] = digits[u&0xF]
		u >>= 4
	}
	i--
	buf[i] = digits[u]
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
