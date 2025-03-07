package xprint

import (
	"reflect"
)

// doPrintf is the core printf implementation. It formats into p.buf.
func (p *printer) printf(format string, args []any) {
	end := len(format)
	argNum := 0
	lenOfArgs := len(args)
	i := 0
	for i < end {
		lasti := i

		for i < end && format[i] != '%' {
			i++
		}

		// if i > lasti && format[i-1] != '%' {
		// 	p.buf = append(p.buf, format[lasti:i]...)
		// }
		switch i > lasti && format[i-1] != '%' {
		case true:
			p.buf = append(p.buf, format[lasti:i]...)
		}

		// if i >= end {
		// 	break
		// }
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

		// Handle flags
		for i < end {
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
		p.arg = args[argNum]
		if i < end && format[i] == '*' {
			i++
			if argNum >= lenOfArgs {
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
			argNum++
		} else if i < end {
			p.fmt.wid, p.fmt.widPresent, i = parsenum(format, i, end)
		}
		// Handle precision
		if i < end && format[i] == '.' {
			i++
			if i < end && format[i] == '*' {
				i++
				if argNum >= lenOfArgs {
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
				argNum++
			} else {
				p.fmt.prec, p.fmt.precPresent, i = parsenum(format, i, end)
			}
		}

		if i >= end {
			p.buf.writeString(noVerbString)
			break
		}

		verb := rune(format[i])
		i++

		// Handle argument
		if argNum >= lenOfArgs {
			p.buf.writeString(missingString)
			break
		}

		argNum++

		if p.arg == nil {
			p.buf.writeString(nilAngleString)
			break
		}
		if p.ArgIsString() && verb == 's' && verb != 'T' {
			p.buf = append(p.buf, p.arg.(string)...)
			continue
		} else if p.ArgIsBytes() && verb == 's' && verb != 'T' {
			p.buf = append(p.buf, p.arg.([]byte)...)
			continue
		}

		switch verb {
		case 'v':
			p.fmt.plusV = p.fmt.plus
			p.fmt.sharpV = p.fmt.sharp
			p.printArg(p.arg, verb)
		case 'd', 'o', 'O', 'x', 'X', 'b', 'B':
			p.printArg(p.arg, verb)
		case 'f', 'F', 'g', 'G', 'e', 'E':
			p.printArg(p.arg, verb)
		case 'q', 's': // 's'
			p.printArg(p.arg, verb)
		case 't':
			p.printBool(p.arg)

		case 'T':
			p.printReflectType(p.arg)
		case 'p':
			p.fmtPointer(reflect.ValueOf(p.arg), verb)
		default:
			p.buf.writeString(percentBangString)
			p.buf.writeRune(verb)
			p.buf.writeString(noVerbString)
		}

	}
}
