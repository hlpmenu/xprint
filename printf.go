package xprint

import (
	"strings"
)

// onlyContainsStringPlaceholders checks if the format string only contains %s placeholders
// and not any other verbs or flags
func onlyContainsStringPlaceholders(format string) bool {
	verbFound := false
	for i := 0; i < len(format); {
		// Fast path for non-% characters
		if format[i] != '%' {
			i++
			continue
		}

		// We have a '%', check next character
		if i+1 >= len(format) {
			// Trailing % with no verb, not a valid format string
			return false
		}

		// Handle escaped percent: %%
		if format[i+1] == '%' {
			i += 2 // Skip both % characters
			continue
		}

		// Process potential verb
		i++ // Skip past the %

		// Skip any digits for argument position if present (for indexed formatting)
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			i++
		}

		// If we reach the end without a verb, not valid
		if i >= len(format) {
			return false
		}

		// Check for flags/width/precision which we don't support in fast path
		// (anything that's not the 's' verb directly)
		switch format[i] {
		case '.', '+', '-', ' ', '#', '0', '*':
			// Flags or width/precision indicators not supported in fast path
			return false
		case 's':
			// Found %s, continue checking
			verbFound = true
			i++
		default:
			// Not %s, not suitable for fast path
			return false
		}
	}

	// Need at least one %s verb to use the fast path
	return verbFound
}

// allArgsAreStringLike checks if all arguments are strings or []byte
func allArgsAreStringLike(args []any) bool {
	// Early return for empty args
	if len(args) == 0 {
		return false
	}

	// Check each argument type
	for _, arg := range args {
		if arg == nil {
			return false
		}

		// Use type assertions directly for better performance
		_, isString := arg.(string)
		_, isByteSlice := arg.([]byte)

		if !isString && !isByteSlice {
			return false
		}
	}
	return true
}

// fastStringFormat is a fast implementation for format strings that only contain %s
// and where all arguments are strings or []byte
func fastStringFormat(format string, args []any) string {
	// Preallocate a string builder
	var result strings.Builder

	argIndex := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			result.WriteByte(format[i])
			continue
		}

		// Handle escaped percent: %%
		if i+1 < len(format) && format[i+1] == '%' {
			result.WriteByte('%')
			i++
			continue
		}

		// Process a %s verb
		i++

		// Skip any digits for argument position if present
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			i++
		}

		// We have a %s placeholder, replace it with the argument
		if argIndex < len(args) {
			switch arg := args[argIndex].(type) {
			case string:
				result.WriteString(arg)
			case []byte:
				result.Write(arg)
			}
			argIndex++
		}
	}

	return result.String()
}
