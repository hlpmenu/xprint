package xprint

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
