package xprint

import (
	"reflect"
	"strconv"
)

func (p *printer) print(args []any) {
	var lastWasString bool
	for i, arg := range args {
		if i > 0 && (!lastWasString || !printisNumeric(arg)) {
			p.buf.writeByte(' ')
		}

		switch v := arg.(type) {
		case nil:
			p.buf.writeString(nilAngleString)
			lastWasString = false
		case string:
			p.buf.writeString(v)
			lastWasString = true
		case []byte:
			p.buf.writeByte('[')
			for i, b := range v {
				if i > 0 {
					p.buf.writeByte(' ')
				}
				p.buf.writeString(strconv.Itoa(int(b)))
			}
			p.buf.writeByte(']')
			lastWasString = false
		case bool:
			if v {
				p.buf.writeString("true")
			} else {
				p.buf.writeString("false")
			}
			lastWasString = false
		case int, int8, int16, int32, int64:
			p.fmt.uintbase = 10
			p.printInt(v, 'v')
			lastWasString = false
		case uint, uint8, uint16, uint32, uint64, uintptr:
			p.fmt.uintbase = 10
			p.printInt(v, 'v')
			lastWasString = false
		case float32:
			p.printFloat32(v, 'v')
			lastWasString = false
		case float64:
			p.printFloat64(v, 'v')
			lastWasString = false
		case complex64, complex128:
			p.printComplex(v, 'v')
			lastWasString = false
		case error:
			p.buf.writeString(v.Error())
			lastWasString = false
		default:
			// For any other type, use reflection
			p.value = reflect.ValueOf(v)
			if p.value.Kind() == reflect.Ptr {
				if p.value.IsNil() {
					p.buf.writeString(nilAngleString)
				} else {
					p.buf.writeString("0x")
					p.buf.writeString(strconv.FormatUint(uint64(p.value.Pointer()), 16))
				}
			} else {
				p.printValue(p.value, 'v', 0)
			}
			lastWasString = false
		}
	}
}

func printisNumeric(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return true
	default:
		return false
	}
}
