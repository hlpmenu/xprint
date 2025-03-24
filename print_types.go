package xprint

import (
	"reflect"
	"unsafe"
)

func (p *printer) fmtPointer(value any, verb rune) {
	// For nil input, return immediately
	if value == nil {
		p.buf.writeString(nilAngleString)
		return
	}

	v := reflect.ValueOf(value)

	// Handle nil pointers
	if v.Kind() == reflect.Pointer && v.IsNil() {
		p.buf.writeString(nilAngleString)
		return
	}

	// For 'p' verb, always use hexadecimal address
	if verb == 'p' {
		var u uintptr

		switch x := value.(type) {
		case unsafe.Pointer:
			u = uintptr(x)
		case uintptr:
			u = x
		case reflect.Value:
			if x.Kind() == reflect.Pointer {
				u = x.Pointer()
			} else {
				p.buf.writeString(nilAngleString)
				return
			}
		default:
			if v.Kind() == reflect.Pointer {
				u = v.Pointer()
			} else {
				p.buf.writeString(nilAngleString)
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
		return
	}

	// If verb is 'v', use special handling for common types
	if verb == 'v' && v.Kind() == reflect.Pointer {
		// If it's a pointer to a struct or a common type, dereference it
		elem := v.Elem()
		if elem.IsValid() {
			p.printValue(elem, verb, 0)
			return
		}
	}

	// Default to hex pointer
	var u uintptr

	switch x := value.(type) {
	case unsafe.Pointer:
		u = uintptr(x)
	case uintptr:
		u = x
	case reflect.Value:
		if x.Kind() == reflect.Pointer {
			u = x.Pointer()
		} else {
			p.buf.writeString(nilAngleString)
			return
		}
	default:
		if v.Kind() == reflect.Pointer {
			u = v.Pointer()
		} else {
			p.buf.writeString(nilAngleString)
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
		// For the 'v' verb, use 'g' format for both real and imaginary parts
		realVerb := verb
		if verb == 'v' {
			realVerb = 'g'
		}
		p.printFloat32(real(v), realVerb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat32(imag(v), realVerb)
	case complex128:
		// For the 'v' verb, use 'g' format for both real and imaginary parts
		realVerb := verb
		if verb == 'v' {
			realVerb = 'g'
		}
		p.printFloat64(real(v), realVerb)
		if imag(v) >= 0 {
			p.buf.writeByte('+')
		}
		p.printFloat64(imag(v), realVerb)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	p.buf.writeByte('i')
	p.buf.writeByte(')')
}
