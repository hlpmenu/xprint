package xprint

import (
	"strconv"
)

func (p *printer) printInt(v any, verb rune) {
	var str string
	switch v := v.(type) {
	case int:
		str = p.printIntFast(v)
	case int8:
		str = p.printInt8Fast(v)
	case int16:
		str = p.printInt16Fast(v)
	case int32:
		str = p.printInt32Fast(v)
	case int64:
		str = p.printInt64Fast(v)
	case uint:
		str = strconv.FormatUint(uint64(v), p.fmt.uintbase)
	case uint8:
		str = strconv.FormatUint(uint64(v), p.fmt.uintbase)
	case uint16:
		str = strconv.FormatUint(uint64(v), p.fmt.uintbase)
	case uint32:
		str = strconv.FormatUint(uint64(v), p.fmt.uintbase)
	case uint64:
		str = strconv.FormatUint(v, p.fmt.uintbase)
	case uintptr:
		str = strconv.FormatUint(uint64(v), p.fmt.uintbase)
	default:
		p.buf.writeString(percentBangString)
		p.buf.writeByte(byte(verb))
		p.buf.writeString(badVerbString)
		return
	}
	switch p.fmt.toupper {
	case true:
		p.buf.writeStringToUpper(str)
	case false:
		p.buf.writeString(str)
	}
}

func (p *printer) printIntFast(v int) string {
	switch v {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return strconv.FormatInt(int64(v), p.fmt.uintbase)
	}
}

func (p *printer) printInt8Fast(v int8) string {
	switch v {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return strconv.FormatInt(int64(v), p.fmt.uintbase)
	}
}

func (p *printer) printInt16Fast(v int16) string {
	switch v {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return strconv.FormatInt(int64(v), p.fmt.uintbase)
	}
}

func (p *printer) printInt32Fast(v int32) string {
	switch v {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return strconv.FormatInt(int64(v), p.fmt.uintbase)
	}
}

func (p *printer) printInt64Fast(v int64) string {
	switch v {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return strconv.FormatInt(v, p.fmt.uintbase)
	}
}

type num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func predefinedOrPrint[T num](v T, base int, verb rune) string {
	var str string
	switch v {
	case 0:
		str = "0"
	case 1:
		str = "1"
	case 2:
		str = "2"
	case 3:
		str = "3"
	case 4:
		str = "4"
	case 5:
		str = "5"
	case 6:
		str = "6"
	case 7:
		str = "7"
	case 8:
		str = "8"
	case 9:
		str = "9"
	default:
		str = strconv.FormatInt(int64(v), base)
	}
	return str
}
