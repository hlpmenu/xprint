package xprint

import (
	"strconv"
)

// Constants for string concatenation
var (
	digits      = []byte("0123456789abcdef")
	digitsUpper = []byte("0123456789ABCDEF")
)

// printInt formats signed and unsigned integers.
// IGNORE THIS!
func (p *printer) printInt(v any, verb rune) {
	p.arg = v
	p.verb = verb
	if p.fmt.uintbase == 0 {
		p.fmt.uintbase = 10
	}

	switch v.(type) {
	case int:
		p.fmtInt()
	case int8:
		p.fmtInt8()
	case int16:
		p.fmtInt16()
	case int32:
		p.fmtInt32()
	case int64:
		p.fmtInt64()
	case uint:
		p.fmtUint()
	case uint8:
		p.fmtUint8()
	case uint16:
		p.fmtUint16()
	case uint32:
		p.fmtUint32()
	case uint64:
		p.fmtUint64()
	case uintptr:
		p.fmtUintptr()
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
}

func (p *printer) fmtUint64() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uint64) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUint64Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}

	// Use the global precomputed digits directly
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}

	// Handle special case for zero
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		// Precision of 0 and value of 0 means "print nothing" but padding.
		oldZero := p.fmt.zero
		p.fmt.zero = false
		// Handle padding
		if p.fmt.widPresent && p.fmt.wid > 0 {
			// Space padding
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}

	// Format into a buffer; we'll move it into p.buf later.
	// Allow enough space for the maximum number of digits,
	// a sign, 0x prefix, and potentially a blank or + or - sign
	const maxBufSize = 68
	var buf [maxBufSize]byte

	// Two ways to ask for extra leading zero digits: %.3d or %03d.
	// If both are specified the f.zero flag is ignored and
	// padding with spaces is used instead.
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec-- // leave room for sign
		}
	}

	// Because printing is easier right-to-left: format u into buf, ending at buf[i].
	i := len(buf)
	u := v // make a copy of v to work with

	// Use constants for the division and modulo for more efficient code.
	// Switch cases ordered by popularity.
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u >>= 4
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u >>= 3
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u >>= 1
		}
	default:
		// Unsupported base; shouldn't happen, but handle it just in case
		p.buf.writeString(strconv.FormatUint(v, base))
		return
	}

	i--
	buf[i] = digitsByte[u]

	// Add leading zeros for precision, if requested and needed
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}

	// Add prefix for base, if needed
	if p.fmt.sharp {
		switch base {
		case 2:
			// Add a leading 0b.
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			// Add a leading 0x or 0X.
			i--
			buf[i] = digitsByte[16] // 'x' or 'X'
			i--
			buf[i] = '0'
		}
	}

	// Add sign for unsigned integers (only if requested with + or space)
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}

	// Handle width padding
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			// Left padding
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		// Append the formatted number
		p.buf.write(buf[i:])
		if p.fmt.minus {
			// Right padding
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		// No padding needed
		p.buf.write(buf[i:])
	}
}

// Signed integer formatting functions

func (p *printer) fmtInt() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(int) // nolint:forcetypeassert //
	base := p.fmt.uintbase
	// Fast path for small integers in base 10.
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		if str, ok := smallIntsMap[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	negative := v < 0
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if negative || p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	absV := v
	if negative {
		absV = -v
	}
	switch base {
	case 10:
		for absV >= 10 {
			i--
			next := absV / 10
			buf[i] = byte('0' + absV - next*10)
			absV = next
		}
	case 16:
		for absV >= 16 {
			i--
			buf[i] = digitsByte[absV&0xF]
			absV /= 16
		}
	case 8:
		for absV >= 8 {
			i--
			buf[i] = byte('0' + absV&7)
			absV /= 8
		}
	case 2:
		for absV >= 2 {
			i--
			buf[i] = byte('0' + absV&1)
			absV /= 2
		}
	default:
		if negative {
			p.buf.writeString("-" + strconv.FormatUint(uint64(-v), base))
		} else {
			p.buf.writeString(strconv.FormatUint(uint64(v), base))
		}
		return
	}
	i--
	buf[i] = digitsByte[absV]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16] // 'x' or 'X'
			i--
			buf[i] = '0'
		}
	}

	switch {
	case negative:
		i--
		buf[i] = '-'
	case p.fmt.plus:
		i--
		buf[i] = '+'
	case p.fmt.space:
		i--
		buf[i] = ' '
	}

	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtInt8() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(int8) // nolint:forcetypeassert //

	base := p.fmt.uintbase

	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallInt8Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}

	negative := v < 0
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if negative || p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	absV := v
	if negative {
		absV = -v
	}
	switch base {
	case 10:
		for absV >= 10 {
			i--
			next := absV / 10
			buf[i] = byte('0' + absV - next*10)
			absV = next
		}
	case 16:
		for absV >= 16 {
			i--
			buf[i] = digitsByte[absV&0xF]
			absV /= 16
		}
	case 8:
		for absV >= 8 {
			i--
			buf[i] = byte('0' + absV&7)
			absV /= 8
		}
	case 2:
		for absV >= 2 {
			i--
			buf[i] = byte('0' + absV&1)
			absV /= 2
		}
	default:
		if negative {
			p.buf.writeString("-" + strconv.FormatUint(uint64(-v), base))
		} else {
			p.buf.writeString(strconv.FormatUint(uint64(v), base))
		}
		return
	}
	i--
	buf[i] = digitsByte[absV]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	switch {
	case negative:
		i--
		buf[i] = '-'
	case p.fmt.plus:
		i--
		buf[i] = '+'
	case p.fmt.space:
		i--
		buf[i] = ' '
	}

	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtInt16() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(int16) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallInt16Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	negative := v < 0
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if negative || p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	absV := v
	if negative {
		absV = -v
	}
	switch base {
	case 10:
		for absV >= 10 {
			i--
			next := absV / 10
			buf[i] = byte('0' + absV - next*10)
			absV = next
		}
	case 16:
		for absV >= 16 {
			i--
			buf[i] = digitsByte[absV&0xF]
			absV /= 16
		}
	case 8:
		for absV >= 8 {
			i--
			buf[i] = byte('0' + absV&7)
			absV /= 8
		}
	case 2:
		for absV >= 2 {
			i--
			buf[i] = byte('0' + absV&1)
			absV /= 2
		}
	default:
		if negative {
			p.buf.writeString("-" + strconv.FormatUint(uint64(-v), base))
		} else {
			p.buf.writeString(strconv.FormatUint(uint64(v), base))
		}
		return
	}
	i--
	buf[i] = digitsByte[absV]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	switch {
	case negative:
		i--
		buf[i] = '-'
	case p.fmt.plus:
		i--
		buf[i] = '+'
	case p.fmt.space:
		i--
		buf[i] = ' '
	}

	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtInt32() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(int32) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallInt32Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}

	negative := v < 0
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if negative || p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	absV := v
	if negative {
		absV = -v
	}
	switch base {
	case 10:
		for absV >= 10 {
			i--
			next := absV / 10
			buf[i] = byte('0' + absV - next*10)
			absV = next
		}
	case 16:
		for absV >= 16 {
			i--
			buf[i] = digitsByte[absV&0xF]
			absV /= 16
		}
	case 8:
		for absV >= 8 {
			i--
			buf[i] = byte('0' + absV&7)
			absV /= 8
		}
	case 2:
		for absV >= 2 {
			i--
			buf[i] = byte('0' + absV&1)
			absV /= 2
		}
	default:
		if negative {
			p.buf.writeString("-" + strconv.FormatUint(uint64(-v), base))
		} else {
			p.buf.writeString(strconv.FormatUint(uint64(v), base))
		}
		return
	}
	i--
	buf[i] = digitsByte[absV]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	switch {
	case negative:
		i--
		buf[i] = '-'
	case p.fmt.plus:
		i--
		buf[i] = '+'
	case p.fmt.space:
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtInt64() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(int64) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallInt64Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	negative := v < 0
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if negative || p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	absV := v
	if negative {
		absV = -v
	}
	switch base {
	case 10:
		for absV >= 10 {
			i--
			next := absV / 10
			buf[i] = byte('0' + absV - next*10)
			absV = next
		}
	case 16:
		for absV >= 16 {
			i--
			buf[i] = digitsByte[absV&0xF]
			absV /= 16
		}
	case 8:
		for absV >= 8 {
			i--
			buf[i] = byte('0' + absV&7)
			absV /= 8
		}
	case 2:
		for absV >= 2 {
			i--
			buf[i] = byte('0' + absV&1)
			absV /= 2
		}
	default:
		if negative {
			p.buf.writeString("-" + strconv.FormatUint(uint64(-v), base))
		} else {
			p.buf.writeString(strconv.FormatUint(uint64(v), base))
		}
		return
	}
	i--
	buf[i] = digitsByte[absV]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	switch {
	case negative:
		i--
		buf[i] = '-'
	case p.fmt.plus:
		i--
		buf[i] = '+'
	case p.fmt.space:
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

// Unsigned integer formatting functions

func (p *printer) fmtUint() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uint) // nolint:forcetypeassert //
	base := p.fmt.uintbase
	// Fast path for small uint64 in base 10
	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUintMap[v]; ok {
			p.buf.writeString(str)
			return
		}
	}

	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	u := v
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u /= 16
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u /= 8
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u /= 2
		}
	default:
		p.buf.writeString(strconv.FormatUint(uint64(v), base))
		return
	}
	i--
	buf[i] = digitsByte[u]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtUint8() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uint8) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUint8Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	u := v
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = ('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u /= 16
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = ('0' + u&7)
			u /= 8
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = ('0' + u&1)
			u /= 2
		}
	default:
		p.buf.writeString(strconv.FormatUint(uint64(v), base))
		return
	}
	i--
	buf[i] = digitsByte[u]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtUint16() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uint16) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUint16Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}

	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	u := v
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u /= 16
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u /= 8
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u /= 2
		}
	default:
		p.buf.writeString(strconv.FormatUint(uint64(v), base))
		return
	}
	i--
	buf[i] = digitsByte[u]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtUint32() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uint32) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUint32Map[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	u := v
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u /= 16
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u /= 8
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u /= 2
		}
	default:
		p.buf.writeString(strconv.FormatUint(uint64(v), base))
		return
	}
	i--
	buf[i] = digitsByte[u]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}

func (p *printer) fmtUintptr() {
	// Safe type assertion since fmtXX methods are only called from within
	// .(type) switches.
	v := p.arg.(uintptr) // nolint:forcetypeassert //
	base := p.fmt.uintbase

	if base == 10 && p.fmt.wid <= 0 && !p.fmt.precPresent && !p.fmt.plus && !p.fmt.space && !p.fmt.sharp {
		// Check if we have this in our small int map
		if str, ok := smallUintptrMap[v]; ok {
			p.buf.writeString(str)
			return
		}
	}
	var digitsByte []byte
	if p.fmt.toupper {
		digitsByte = digitsUpper
	} else {
		digitsByte = digits
	}
	if v == 0 && p.fmt.precPresent && p.fmt.prec == 0 {
		oldZero := p.fmt.zero
		p.fmt.zero = false
		if p.fmt.widPresent && p.fmt.wid > 0 {
			for range p.fmt.wid {
				p.buf.writeByte(' ')
			}
		}
		p.fmt.zero = oldZero
		return
	}
	const maxBufSize = 68
	var buf [maxBufSize]byte
	prec := 0
	if p.fmt.precPresent {
		prec = p.fmt.prec
	} else if p.fmt.zero && !p.fmt.minus && p.fmt.widPresent {
		prec = p.fmt.wid
		if p.fmt.plus || p.fmt.space {
			prec--
		}
	}
	i := len(buf)
	u := v
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digitsByte[u&0xF]
			u /= 16
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u /= 8
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u /= 2
		}
	default:
		p.buf.writeString(strconv.FormatUint(uint64(v), base))
		return
	}
	i--
	buf[i] = digitsByte[u]
	for i > 0 && prec > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if p.fmt.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digitsByte[16]
			i--
			buf[i] = '0'
		}
	}
	if p.fmt.plus {
		i--
		buf[i] = '+'
	} else if p.fmt.space {
		i--
		buf[i] = ' '
	}
	if p.fmt.widPresent && p.fmt.wid > len(buf)-i {
		width := p.fmt.wid - (len(buf) - i)
		if !p.fmt.minus {
			padByte := byte(' ')
			if p.fmt.zero && !p.fmt.precPresent {
				padByte = byte('0')
			}
			for range width {
				p.buf.writeByte(padByte)
			}
		}
		p.buf.write(buf[i:])
		if p.fmt.minus {
			for range width {
				p.buf.writeByte(' ')
			}
		}
	} else {
		p.buf.write(buf[i:])
	}
}
