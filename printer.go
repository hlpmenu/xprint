package xprint

import (
	"strconv"
	"sync"

	reflect "github.com/goccy/go-reflect"
)

// pp is used to store a printer's state
type printer struct {
	// Big fields first
	buf         buffer
	value       reflect.Value
	arg         any
	visitedPtrs visited
	// wrappedErrs []int
	fmt fmt

	// Frequently updated small fields
	argNum int
	verb   rune

	// Grouped booleans
	recursing  bool
	reordered  bool //nolint:unused
	goodArgNum bool //nolint:unused
	panicking  bool //nolint:unused
	erroring   bool //nolint:unused
	wrapErrs   bool //nolint:unused
}

// func (p *printer) argAsString() string {
// 	return p.arg.(string)
// }

func (p *printer) ArgIsString() bool {
	_, ok := p.arg.(string)
	return ok
}

func (p *printer) ArgIsBytes() bool {
	_, ok := p.arg.([]byte)
	return ok
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *printer {
	// We know this is safe because we're using a sync.Pool
	p := ppFree.Get().(*printer) //nolint:forcetypeassert,errcheck //
	p.fmt.init(&p.buf)
	p.visitedPtrs.init()
	p.recursing = false
	return p
}

// free saves used pp structs in ppFree; avoids an allocation per invocation.
func (p *printer) free() {
	if cap(p.buf) > 1024*64 {
		p.buf = nil
	} else {
		p.buf = p.buf[:0]
	}
	p.arg = nil
	p.value = reflect.Value{}
	p.visitedPtrs.ptrs = nil
	p.recursing = false
	ppFree.Put(p)
}

var ppFree = &sync.Pool{
	New: func() any { return new(printer) },
}

// argNumber returns the next argument to evaluate, which is either the value of the passed-in
// argNum or the value of the bracketed integer that begins format[i:]. It also returns
// the new value of i, that is, the index of the next byte of the format to process.
// func (p *printer) argNumber(argNum int, format string, i int, numArgs int) (newArgNum, newi int, found bool) {
// 	if len(format) <= i || format[i] != '[' {
// 		return argNum, i, false
// 	}
// 	p.reordered = true
// 	index, wid, ok := parseArgNumber(format[i:])
// 	if ok && 0 <= index && index < numArgs {
// 		return index, i + wid, true
// 	}
// 	p.goodArgNum = false
// 	return argNum, i + wid, ok
// }

// func (p *printer) badArgNum(verb rune) {
// 	p.buf.writeString(percentBangString)
// 	p.buf.writeRune(verb)
// 	p.buf.writeString(badIndexString)
// }

// func (p *printer) missingArg(verb rune) {
// 	p.buf.writeString(percentBangString)
// 	p.buf.writeRune(verb)
// 	p.buf.writeString(missingString)
// }

// func (p *printer) writeStringArg() {
// 	p.buf = append(p.buf, p.arg.(string)...)

// }

func (p *printer) printBadVerb(verb rune) {
	p.buf.writeString(percentBangString)
	p.buf.writeRune(verb)
	p.buf.writeString(badVerbString)
}

func (p *printer) printFloat(v any, verb rune) {
	var str string
	switch v := v.(type) {
	case float32:
		str = strconv.FormatFloat(float64(v), byte(verb), p.fmt.prec, 32)
	case float64:
		str = strconv.FormatFloat(v, byte(verb), p.fmt.prec, 64)
	default:
		p.printBadVerb(verb)
		return
	}
	p.buf.writeString(str)
}

// handleMethods checks if the argument implements special formatting interfaces.
func (p *printer) catchPanic(arg any, verb rune, method string) {
	if err := recover(); err != nil {
		// If it's a nil pointer, just say "<nil>"
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Ptr && v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}

		// Otherwise print a concise panic message
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString("(PANIC=")
		p.buf.writeString(method)
		p.buf.writeString(" method: ")
		p.arg = err
		p.verb = 'v'
		p.printArg()
		p.buf.writeByte(')')
	}
}

// func (p *printer) writeNilArg(verb rune) {
// 	p.buf.writeString(percentBangString)
// 	p.buf.writeRune(verb)
// 	p.buf.writeString(nilParenString)
// }
