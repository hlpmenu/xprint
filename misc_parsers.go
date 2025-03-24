package xprint

// parseArgNumber returns the value of the bracketed number, minus 1
// (explicit argument numbers are one-indexed but we want zero-indexed).
// The opening bracket is known to be present at format[0].
// The returned values are the index, the number of bytes to consume
// up to the closing paren, if present, and whether the number parsed
// ok. The bytes to consume will be 1 if no closing paren is present.
// func parseArgNumber(format string) (int, int, bool) {
// 	// There must be at least 3 bytes: [n].
// 	if len(format) < 3 {
// 		return 0, 1, false
// 	}

// 	// Find closing bracket.
// 	for i := 1; i < len(format); i++ {
// 		if format[i] == ']' {
// 			width, ok, newi := parsenum(format, 1, i)
// 			if !ok || newi != i {
// 				return 0, i + 1, false
// 			}
// 			return width - 1, i + 1, true // arg numbers are one-indexed and skip paren.
// 		}
// 	}
// 	return 0, 1, false
// }

func parsenum(s string, start, end int) (int, bool, int) {
	if start >= end {
		return 0, false, end
	}
	num := 0
	isnum := false
	newi := start
	for ; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return num, isnum, newi
}
