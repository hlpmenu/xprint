package xprint

import (
	"reflect"
	"strconv"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

// printArg formats arg in the manner specified by the verb
// and appends it to p.buf.
func (p *printer) printArg(arg any, verb rune) {

	// Handle nil
	if arg == nil {
		logger.Warn("Trigger: arg == nil")
		switch verb {
		case 'T', 'v':
			p.buf.writeString(nilString)
		default:
			p.buf.writeString(percentBangString)
			p.buf = append(p.buf, byte(verb))
			p.buf.writeString(nilString)
		}
		return
	}

	// Handle based on type and verb
	logger.Logf("verb: %s", string(verb))
	switch verb {
	case 'T':
		p.printReflectType(arg)
		return
	case 't':
		p.printBool(arg)
		return
	case 'p':
		logger.Printf("Trigger: p in printarg")
		p.fmtPointer(reflect.ValueOf(arg), verb)
	}
	// Handle by type
	switch v := arg.(type) {
	case []byte:
		p.buf = append(*p.fmt.buf, v...)
	case string:
		p.buf = append(p.buf, v...)
	case bool:
		switch verb {
		case 't':
			p.printBool(v)
		case 's':
			boolstr := percentBangString + "s(" + "bool" + "=" + strconv.FormatBool(v) + ")"
			p.buf = append(p.buf, boolstr...)
		}
	case int, int8, int16, int32, int64:
		p.printInt(v, 10, verb)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		p.printInt(v, 10, verb)
	case float32, float64:
		p.printFloat(v, verb)
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
