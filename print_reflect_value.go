package xprint

import "reflect"

// printValue is similar to printArg but starts with a reflect value, not an interface{} value.
func (p *printer) printValue(v reflect.Value, verb rune, prec int) {
	// Handle nil
	if !v.IsValid() {
		p.buf.writeString(nilAngleString)
		return
	}

	// Check for recursive pointer/interface values
	if !p.recursing && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		ptr := v.Pointer()
		if ptr != 0 && p.visitedPtrs.visit(ptr) {
			// Already seen this pointer, print type and address
			p.buf.writeByte('&')
			p.buf.writeString(v.Type().String())
			p.buf.writeString("(CYCLIC REFERENCE)")
			return
		}
	}

	// Handle special cases for verb 'v' with sharp flag
	if verb == 'v' && p.fmt.sharpV {
		// Print type for nil pointer/interface/slice
		if (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface || v.Kind() == reflect.Slice) && v.IsNil() {
			p.buf.writeString(v.Type().String())
			p.buf.writeString(nilParenString)
			return
		}
		// Print type for other values
		p.buf.writeString(v.Type().String())
		if v.Kind() == reflect.Struct {
			p.buf.writeByte('{')
		} else {
			p.buf.writeByte('(')
		}
	}

	// Set recursing flag for nested calls
	wasRecursing := p.recursing
	p.recursing = true
	defer func() { p.recursing = wasRecursing }()

	// Handle common types
	switch v.Kind() {
	case reflect.Bool:
		p.fmt.fmtBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.printInt(v.Int(), 10, verb)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.printInt(v.Uint(), 10, verb)
	case reflect.Float32, reflect.Float64:
		p.printFloat(v, verb)
	case reflect.String:
		p.fmt.fmtString(v.String())
	case reflect.Slice:
		if v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}
		if v.Type().Elem().Kind() == reflect.Uint8 {
			p.fmt.fmtBytes(v.Bytes())
		} else {
			p.buf.writeByte('[')
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					p.buf.writeByte(' ')
				}
				p.printValue(v.Index(i), verb, prec)
			}
			p.buf.writeByte(']')
		}
	case reflect.Array:
		p.buf.writeByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				p.buf.writeByte(' ')
			}
			p.printValue(v.Index(i), verb, prec)
		}
		p.buf.writeByte(']')
	case reflect.Map:
		if v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}
		p.buf.writeString("map[")
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				p.buf.writeByte(' ')
			}
			p.printValue(key, verb, prec)
			p.buf.writeByte(':')
			p.printValue(v.MapIndex(key), verb, prec)
		}
		p.buf.writeByte(']')
	case reflect.Struct:
		p.buf.writeByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				p.buf.writeByte(' ')
			}
			if p.fmt.plusV {
				p.buf.writeString(v.Type().Field(i).Name)
				p.buf.writeByte(':')
			}
			p.printValue(v.Field(i), verb, prec)
		}
		p.buf.writeByte('}')
	case reflect.Pointer:
		if v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}
		p.buf.writeByte('&')
		p.printValue(v.Elem(), verb, prec)
	case reflect.Interface:
		if v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}
		p.printValue(v.Elem(), verb, prec)
	default:
		// For other types, just use String()
		if v.CanInterface() {
			p.printArg(v.Interface(), verb)
		} else {
			p.buf.writeString(v.String())
		}
	}

	// Close type wrapper for verb 'v' with sharp flag
	if verb == 'v' && p.fmt.sharpV {
		if v.Kind() == reflect.Struct {
			p.buf.writeByte('}')
		} else {
			p.buf.writeByte(')')
		}
	}
}
