package xprint

import (
	"reflect"
	"unsafe"
)

func (p *printer) fmtPointer(value any, verb rune) {
	var u uintptr
	switch v := value.(type) {
	case unsafe.Pointer:
		u = uintptr(v)
	case uintptr:
		u = v
	case reflect.Value:
		u = v.Pointer()
	default:
		switch verb {
		case 's', 'p', 'v':
			// Do nothing
		default:
			p.buf.writeString(nilParenString)
			return
		}
	}

	p.buf.writeByte('0')
	p.buf.writeByte('x')

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

func (p *printer) printBool(arg any) {
	if b, ok := arg.(bool); ok {
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
		p.printFloat(real(v), verb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat(imag(v), verb)
	case complex128:
		p.printFloat(real(v), verb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat(imag(v), verb)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	p.buf.writeByte('i')
	p.buf.writeByte(')')
}
