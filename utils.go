package xprint

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	const maxLimit int = 1e6
	return x > maxLimit || x < -maxLimit
}
