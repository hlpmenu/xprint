package xprint

import (
	"math"
	"reflect"
	"strconv"
)

func (p *printer) printInt(v any, base int, verb rune) {
	var str string
	switch v := v.(type) {
	case int:
		str = predefinedOrPrint(v, base, verb)
	case int8:
		str = predefinedOrPrint(v, base, verb)
	case int16:
		str = predefinedOrPrint(v, base, verb)
	case int32:
		str = predefinedOrPrint(v, base, verb)
	case int64:
		str = predefinedOrPrint(v, base, verb)

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
				// if int64(n) >= 0 && uint64(int(n)) == n {
				// 	num = int(n)
				// 	isInt = true
				// }
				if n <= math.MaxInt {
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
