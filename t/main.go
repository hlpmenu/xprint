package xprint

import (
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

const (
	commaSpaceString  = ", "
	nilAngleString    = "<nil>"
	nilParenString    = "(nil)"
	nilString         = "nil"
	percentBangString = "%!"
	missingString     = "(MISSING)"
	badIndexString    = "(BADINDEX)"
	noVerbString      = "%!(NOVERB)"
	badWidthString    = "%!(BADWIDTH)"
	badPrecString     = "%!(BADPREC)"
	badVerbString     = "%!(BADVERB)"
	mapString         = "map["
	panicString       = "(PANIC="
	extraString       = "%!(EXTRA "

	invReflectString = "<invalid reflect.Value>"
)

// Digits for formatting
const (
	ldigits = "0123456789abcdef"
	udigits = "0123456789ABCDEF"
)

const (
	signed   = true
	unsigned = false
)

// buffer is a simple []byte buffer for building strings.
type buffer []byte

func (b *buffer) Len() int {
	return len([]byte(*b))
}
func (b *buffer) LenMB() int {
	return BtoMB(b.Len())
}

func (b *buffer) write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}

func (p *pp) writeStringArg() {
	p.buf = append(p.buf, p.arg.(string)...)

}

func BtoMB(b int) int {
	return b / 1024 / 1024
}

func (b *buffer) writeByte(c byte) {
	*b = append(*b, c)
}

func (b *buffer) writeRune(r rune) {
	*b = utf8.AppendRune(*b, r)
}

// fmtFlags contains the core formatting flags
type fmtFlags struct {
	widPresent, precPresent         bool
	minus, plus, sharp, space, zero bool
	plusV, sharpV                   bool
	wid, prec                       int
}

// fmt holds the formatting state
type fmt struct {
	buf *buffer
	fmtFlags
	// intbuf is large enough to store %b of an int64 with a sign and
	// avoids padding at the end of the struct on 32 bit architectures.
	intbuf [68]byte
}

func (f *fmt) init(b *buffer) {
	f.buf = b
	f.clearflags()
}

func (f *fmt) clearflags() {
	f.widPresent = false
	f.precPresent = false
	f.minus = false
	f.plus = false
	f.sharp = false
	f.space = false
	f.zero = false
	f.plusV = false
	f.sharpV = false
	f.wid = 0
	f.prec = 0
}

// visited tracks pointers already seen during recursive value formatting
type visited struct {
	ptrs map[uintptr]bool
}

func (v *visited) init() {
	v.ptrs = make(map[uintptr]bool)
}

func (v *visited) visit(p uintptr) bool {
	if v.ptrs[p] {
		return true
	}
	v.ptrs[p] = true
	return false
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	//p.buf = make([]byte, 0, 1024*1024*17)
	p.fmt.init(&p.buf)
	p.visitedPtrs.init()
	//p.fmt.buf = &p.buf
	p.recursing = false
	return p
}

// Printf formats according to a format specifier and returns the resulting string
func Printf(format string, args ...any) string {
	p := newPrinter()
	//	p.doPrintf(format, args)
	p.OlddoPrintf(format, args)
	s := string(p.buf)
	p.free()
	return s
}

func parsenum(s string, start, end int) (num int, isnum bool, newi int) {
	if start >= end {
		return 0, false, end
	}
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return
}

// Stringer is implemented by any value that has a String method.
type Stringer interface {
	String() string
}

// GoStringer is implemented by any value that has a GoString method.
type GoStringer interface {
	GoString() string
}

// handleMethods checks if the argument implements special formatting interfaces.
func (p *pp) catchPanic(arg any, verb rune, method string) {
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

func (p *pp) handleMethods(verb rune) (handled bool) {
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

// printArg formats arg in the manner specified by the verb
// and appends it to p.buf.
func (p *pp) printArg(arg any, verb rune) {

	// Handle nil
	if arg == nil {
		logger.Warn("Trigger: arg == nil")
		switch verb {
		case 'T', 'v':
			p.buf.writeString(nilString)
		default:
			p.buf.writeString(percentBangString)
			p.buf.writeByte(byte(verb))
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

// printValue is similar to printArg but starts with a reflect value, not an interface{} value.
func (p *pp) printValue(v reflect.Value, verb rune, prec int) {
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
		p.printFloat(v.Float(), verb)
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

// Core formatting methods
func (f *fmt) fmtBool(v bool) {
	if v {
		f.buf.writeString("true")
	} else {
		f.buf.writeString("false")
	}
}

func (f *fmt) fmtString(s string) {
	f.buf.writeString(s)
}

func (f *fmt) fmtBytes(v []byte) {
	f.buf.write(v)
}

func (f *fmt) fmtFloat(v float64, size int, verb rune, prec int) {
	// Handle sign
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	} else if f.plus {
		sign = "+"
	} else if f.space {
		sign = " "
	}

	// Set default precision if not specified
	if prec < 0 {
		switch verb {
		case 'e', 'E', 'f', 'F':
			prec = 6
		case 'g', 'G':
			prec = -1 // Use shortest representation
		default:
			prec = 6
		}
	}

	// Format based on verb
	var format byte
	switch verb {
	case 'f', 'F':
		format = 'f'
	case 'e', 'E':
		format = byte(verb)
	case 'g', 'G':
		format = byte(verb)
	default:
		format = 'g'
	}

	// Convert to string
	num := strconv.FormatFloat(v, format, prec, size)

	// Handle special cases for 'g'/'G' format
	if (verb == 'g' || verb == 'G') && f.sharp && strings.IndexByte(num, '.') < 0 {
		num += "."
	}

	// Combine all parts
	s := sign + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) fmtInt(v int64, base int, verb rune) {
	// Build the string
	var s string
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	} else if f.plus {
		sign = "+"
	} else if f.space {
		sign = " "
	}

	// Convert number to string
	num := strconv.FormatInt(v, base)
	if base == 16 && verb == 'X' {
		num = strings.ToUpper(num)
	}

	// Add prefix if requested
	prefix := ""
	if f.sharp {
		switch base {
		case 2:
			prefix = "0b"
		case 8:
			prefix = "0o"
		case 16:
			prefix = "0x"
		}
	}

	// Combine all parts
	s = sign + prefix + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) fmtUint(v uint64, base int, verb rune) {
	// Build the string
	var s string
	sign := ""
	if f.plus {
		sign = "+"
	} else if f.space {
		sign = " "
	}

	// Convert number to string
	num := strconv.FormatUint(v, base)
	if base == 16 && verb == 'X' {
		num = strings.ToUpper(num)
	}

	// Add prefix if requested
	prefix := ""
	if f.sharp {
		switch base {
		case 2:
			prefix = "0b"
		case 8:
			prefix = "0o"
		case 16:
			prefix = "0x"
		}
	}

	// Combine all parts
	s = sign + prefix + num

	// Handle padding
	f.pad(s)
}

func (f *fmt) pad(s string) {
	if !f.widPresent || f.wid == 0 {
		f.buf.writeString(s)
		return
	}

	width := f.wid - utf8.RuneCountInString(s)
	if width <= 0 {
		f.buf.writeString(s)
		return
	}

	// Left padding
	if !f.minus {
		if f.zero && !strings.ContainsAny(s, ".-+") {
			for i := 0; i < width; i++ {
				f.buf.writeByte('0')
			}
		} else {
			for i := 0; i < width; i++ {
				f.buf.writeByte(' ')
			}
		}
	}

	f.buf.writeString(s)

	// Right padding
	if f.minus {
		for i := 0; i < width; i++ {
			f.buf.writeByte(' ')
		}
	}
}

// ... other supporting methods copied from fmt package and adapted ...

func (p *pp) printInt(v any, base int, verb rune) {
	var str string
	switch v := v.(type) {
	case int:
		str = strconv.FormatInt(int64(v), base)
	case int8:
		str = strconv.FormatInt(int64(v), base)
	case int16:
		str = strconv.FormatInt(int64(v), base)
	case int32:
		str = strconv.FormatInt(int64(v), base)
	case int64:
		str = strconv.FormatInt(v, base)
	case uint:
		str = strconv.FormatUint(uint64(v), base)
	case uint8:
		str = strconv.FormatUint(uint64(v), base)
	case uint16:
		str = strconv.FormatUint(uint64(v), base)
	case uint32:
		str = strconv.FormatUint(uint64(v), base)
	case uint64:
		str = strconv.FormatUint(v, base)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	p.buf.writeString(str)
}

func (p *pp) printFloat(v any, verb rune) {
	var str string
	switch v := v.(type) {
	case float32:
		str = strconv.FormatFloat(float64(v), byte(verb), p.fmt.prec, 32)
	case float64:
		str = strconv.FormatFloat(v, byte(verb), p.fmt.prec, 64)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	p.buf.writeString(str)
}

func (p *pp) printString() {
}

func (p *pp) printBool(arg any) {
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

func (p *pp) printReflectType(arg any) {
	p.buf.writeString(reflect.TypeOf(arg).String())
}

func (p *pp) fmtPointer(value any, verb rune) {
	var u uintptr
	switch v := value.(type) {
	case unsafe.Pointer:
		logger.Warn("Trigger: unsafe.Pointer")
		u = uintptr(v)
	case uintptr:
		logger.Warn("Trigger: uintptr")
		u = v
	case reflect.Value:
		logger.Warn("Trigger: reflect.Value")
		u = v.Pointer()
	default:
		logger.Warnf("Trigger: default with verb: %s", string(verb))
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

func (p *pp) printComplex(v any, verb rune) {
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

// intFromArg gets the argNumth element of a. On return, isInt reports whether the argument has integer type.
func intFromArg(a []any, argNum int) (num int, isInt bool, newArgNum int) {
	newArgNum = argNum
	if argNum < len(a) {
		num, isInt = a[argNum].(int) // Almost always OK.
		if !isInt {
			// Work harder.
			switch v := reflect.ValueOf(a[argNum]); v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				n := v.Int()
				if int64(int(n)) == n {
					num = int(n)
					isInt = true
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				n := v.Uint()
				if int64(n) >= 0 && uint64(int(n)) == n {
					num = int(n)
					isInt = true
				}
			default:
				// Already 0, false.
			}
		}
		newArgNum = argNum + 1
		if tooLarge(num) {
			num = 0
			isInt = false
		}
	}
	return
}

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	const max int = 1e6
	return x > max || x < -max
}
