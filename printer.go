package xprint

import (
	"reflect"
	"strconv"
	"sync"
)

// pp is used to store a printer's state
type printer struct {
	buf   buffer
	arg   any
	value reflect.Value
	fmt   fmt
	// Track recursive pointer formatting
	visitedPtrs visited
	recursing   bool
	// reordered records whether the format string used argument reordering.
	reordered bool //nolint:unused //
	// goodArgNum records whether the most recent reordering directive was valid.
	goodArgNum bool //nolint:unused //
	// panicking is set by catchPanic to avoid infinite panic, recover, panic, ... recursion.
	panicking bool //nolint:unused //
	// erroring is set when printing an error string to guard against calling handleMethods.
	erroring bool //nolint:unused //
	// wrapErrs is set when the format string may contain a %w verb.
	wrapErrs bool //nolint:unused //
	// wrappedErrs records the targets of the %w verb.
	wrappedErrs []int //nolint:unused //
	// argNum tracks the current argument number being processed
	argNum int
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
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Pointer && v.IsNil() {
			p.buf.writeString(nilAngleString)
			return
		}

		// Otherwise print a concise panic message
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString("(PANIC=")
		p.buf.writeString(method)
		p.buf.writeString(" method: ")
		p.printArg(err, 'v')
		p.buf.writeByte(')')
	}
}

func (p *printer) handleMethods(verb rune) bool {
	// Handle error values
	if err, ok := p.arg.(error); ok {
		defer p.catchPanic(p.arg, verb, "Error")
		p.fmt.fmtString(err.Error())
		return true
	}

	// Handle GoStringer for %#v
	if p.fmt.sharpV {
		if stringer, ok := p.arg.(GoStringer); ok {
			defer p.catchPanic(p.arg, verb, "GoString")
			p.fmt.fmtString(stringer.GoString())
			return true
		}
	}

	// Handle Stringer for %v, %s
	if verb == 'v' || verb == 's' {
		if stringer, ok := p.arg.(Stringer); ok {
			defer p.catchPanic(p.arg, verb, "String")
			p.fmt.fmtString(stringer.String())
			return true
		}
	}

	return false
}
