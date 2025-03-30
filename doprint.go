package xprint

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

		// switch end <= i {
		// case true:
		// 	break
		// }

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
			current := format[i]
			var next byte
			if i+1 < end {
				next = format[i+1]
			}

			switch {
			case current == '#' && end >= i+1 && next != 'v':
				p.fmt.sharp = true
			case current == '#' && end >= i+1 && next == 'v':
				p.fmt.sharpV = true
			case current == '0':
				p.fmt.zero = true
			case current == '+' && end >= i+1 && next != 'v':
				p.fmt.plus = true
			case current == '+' && end >= i+1 && next == 'v':
				p.fmt.plusV = true
			case current == '-':
				p.fmt.minus = true
			case current == ' ':
				p.fmt.space = true
			default:
				goto flags_done
			}
			i++
		}
	flags_done:
		if i >= end || p.argNum >= lenOfArgs {
			p.buf.writeString(percentBangString)
			p.buf.writeByte(format[i])
			p.buf.writeString(missingString)
			break
		}
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

		if i >= end {
			p.buf.writeString(noVerbString)
			break
		}

		p.verb = rune(format[i])
		i++

		// Handle argument
		if p.argNum >= lenOfArgs {
			p.buf.writeString(missingString)
			break
		}

		p.argNum++

		// switch {
		// case p.arg == nil:

		// 	p.buf.writeNilArg(p.verb)
		// 	continue
		// }

		if p.ArgIsString() && p.verb == 's' && p.verb != 'T' && !p.fmt.widPresent {
			// Fast path: string value with no width formatting, use direct concatenation
			p.buf = append(p.buf, p.arg.(string)...) //nolint:forcetypeassert //
			continue
		} else if p.ArgIsBytes() && p.verb == 's' && p.verb != 'T' && !p.fmt.widPresent {
			// Fast path: byte slice value with no width formatting, use direct concatenation
			p.buf = append(p.buf, p.arg.([]byte)...) //nolint:forcetypeassert //
			continue
		}
		p.fmt.uintbase = 10
		p.fmt.toupper = false
		switch p.verb {
		case 'v':
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
			switch v := p.arg.(type) { //nolint:all //
			case string:
				p.arg = `"` + v + `"`
			}
			p.printArg()
		case 't':
			p.printBool()
		case 'T':
			p.printReflectType(p.arg)
		case 'p':
			p.fmtPointer(p.arg, p.verb)
		default:
			p.buf.writeString(percentBangString)
			p.buf.writeRune(p.verb)
			p.buf.writeString(noVerbString)

		}
	}
}
