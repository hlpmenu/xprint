package xprint

import (
	"log"
	"reflect"
)

// doPrintf is the core printf implementation. It formats into p.buf.
func (p *printer) printf(format string, args []any) {
	end := len(format)
	p.argNum = 0
	lenOfArgs := len(args)
	i := 0
	for i < end {
		lasti := i

		for i < end && format[i] != '%' {
			i++
		}

		if i > lasti && format[i-1] != '%' {
			p.buf = append(p.buf, format[lasti:i]...)
		}

		switch end <= i {
		case true:
			break
		}

		// Process one verb
		i++

		// Handle %% case
		if i < end && format[i] == '%' {
			p.buf = append(p.buf, '%')
			i++
			continue
		}

		p.fmt.clearflags()
		var c int
		// Handle flags
		for i < end {
			c++
			switch format[i] {
			case '#':
				p.fmt.sharp = true
			case '0':
				p.fmt.zero = true
			case '+':
				p.fmt.plus = true
			case '-':
				p.fmt.minus = true
			case ' ':
				p.fmt.space = true
			default:
				goto flags_done
			}
			i++
		}
	flags_done:
		log.Print("flags done loop iteration")
		if i >= end || p.argNum >= len(args) {
			p.buf.writeString(percentBangString)
			p.buf.writeRune(rune(format[i]))
			p.buf.writeString(missingString)
			break
		}
		log.Printf("argnum: %d", p.argNum)
		p.arg = args[p.argNum]

		if i < end && format[i] == '*' {
			i++
			if p.argNum >= lenOfArgs {
				p.buf.writeString(missingString)
				break
			}
			width, ok := p.arg.(int)
			if !ok {
				p.buf.writeString(badWidthString)
			} else {
				p.fmt.wid = width
				p.fmt.widPresent = true
				if width < 0 {
					p.fmt.minus = true
					p.fmt.wid = -width
				}
			}
			p.argNum++
		} else if i < end {
			p.fmt.wid, p.fmt.widPresent, i = parsenum(format, i, end)
		}
		// Handle precision
		if i < end && format[i] == '.' {
			i++
			if i < end && format[i] == '*' {
				i++
				if p.argNum >= lenOfArgs {
					p.buf.writeString(missingString)
					break
				}
				prec, ok := p.arg.(int)
				if !ok {
					p.buf.writeString(badPrecString)
				} else {
					p.fmt.prec = prec
					p.fmt.precPresent = true
					if prec < 0 {
						p.fmt.precPresent = false
					}
				}
				p.argNum++
			} else {
				p.fmt.prec, p.fmt.precPresent, i = parsenum(format, i, end)
			}
		}

		var lastIteration bool
		if i >= end {
			p.buf.writeString(noVerbString)
			break
		}

		p.verb = rune(format[i])
		i++

		if i == end {
			lastIteration = true
		}

		// Handle argument
		if p.argNum >= lenOfArgs {
			p.buf.writeString(missingString)
			break
		}

		p.argNum++

		switch {
		case p.arg == nil:
			p.buf.writeNilArg(p.verb)
			continue
		}

		if p.ArgIsString() && p.verb == 's' && p.verb != 'T' && !p.fmt.widPresent {
			// Fast path: string value with no width formatting, use direct concatenation
			p.buf = append(p.buf, p.arg.(string)...)
			continue
		} else if p.ArgIsBytes() && p.verb == 's' && p.verb != 'T' && !p.fmt.widPresent {
			// Fast path: byte slice value with no width formatting, use direct concatenation
			p.buf = append(p.buf, p.arg.([]byte)...)
			continue
		}
		p.fmt.uintbase = 10
		p.fmt.toupper = false
		switch p.verb {
		case 'v':
			p.fmt.plusV = p.fmt.plus
			p.fmt.sharpV = p.fmt.sharp
			p.printArg()
		case 'o':
			p.fmt.uintbase = 8
			p.printArg()
		case 'O':
			p.fmt.uintbase = 8
			p.fmt.toupper = true
		case 'd':
			p.printArg()
		case 'x':
			p.fmt.uintbase = 16
			p.printArg()
		case 'X':
			p.fmt.uintbase = 16
			p.fmt.toupper = true
			p.printArg()
		case 'b':
			p.fmt.uintbase = 2
			p.printArg()
		case 'B':
			p.fmt.uintbase = 2
			p.fmt.toupper = true
			p.printArg()
		case 'f', 'F', 'g', 'G', 'e', 'E':
			p.printArg()
		case 's': // 's'
			p.printArg()
		case 'q':
			// use switch even tho single case for more oprimal type conv
			switch v := p.arg.(type) {
			case string:
				p.arg = `"` + v + `"`
			}
			p.printArg()
		case 't':
			p.printBool(p.arg)
		case 'T':
			p.printReflectType(p.arg)
		case 'p':
			p.fmtPointer(reflect.ValueOf(p.arg), p.verb)
		default:
			p.buf.writeString(percentBangString)
			p.buf.writeRune(p.verb)
			p.buf.writeString(noVerbString)
			if lastIteration {
				break
			}
		}
	}
}

func Example(b []byte, a ...any) {
	// b is []byte
	_ = b

	// a is either any.([]strings) or []any.(string)

	var holds any

	// if the compiler has casteed it as any.([]string),
	// this work just fine
	holds = a[0]
	_ = holds

	// if were a switch
	for {
		// but if compiler has casted it as a slice of interfaces,
		// which seems to be almost random, this will panic.
		// here a is type [][]any
		holds = a[0]

	}

}

func callexample() {
	// type []string
	input := []string{"hello", "world"}
	Example([]byte{}, input)
}
