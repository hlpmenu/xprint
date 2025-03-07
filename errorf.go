package xprint

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
	"slices"
	"strings"
)

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
//
// If the format specifier includes a %w verb with an error operand,
// the returned error will implement an Unwrap method returning the operand. If there
// is more than one %w verb, the returned error will implement an Unwrap
// method returning a slice of the operands in the order they appear.
// It is invalid to supply the %w verb with an operand that does not implement
// the error interface. The %w verb is otherwise a synonym for %v.
func Errorf(format string, a ...any) error {
	// Fast path for no arguments - just return the format string as an error
	if len(a) == 0 {
		return errors.New(format)
	}
	// Fast path for common cases:
	// 1. Simple error without wrapping: "%s" with a string argument
	// 2. Simple wrapped error: "%w" with an error argument
	if len(a) == 1 {
		if isSimpleFormat(format) {
			if format == "%s" {
				// Simple string error without wrapping
				switch arg := a[0].(type) {
				case string:
					return errors.New(arg)
				case []byte:
					return errors.New(string(arg))
				case error:
					return errors.New(arg.Error())
				}
			} else if format == "%w" {
				// Simple error wrapping
				if err, ok := a[0].(error); ok {
					return &wrapError{
						msg: err.Error(),
						err: err,
					}
				}
			}
		}
	}

	// Find %w verbs and replace them with %v, while tracking positions
	wrappedErrs := make([]int, 0, 1)
	reordered := false
	modifiedFormat := parseErrorFormat(format, &wrappedErrs, &reordered)

	// Format the message using the modified format
	s := Printf(modifiedFormat, a...)

	// Create appropriate error type
	var err error
	switch len(wrappedErrs) {
	case 0:
		// No wrapped errors, just a regular error
		err = errors.New(s)
	case 1:
		// Single wrapped error
		w := &wrapError{msg: s}
		w.err, _ = a[wrappedErrs[0]].(error)
		err = w
	default:
		// Multiple wrapped errors
		if reordered {
			slices.Sort(wrappedErrs)
		}
		var errs []error
		for i, argNum := range wrappedErrs {
			if i > 0 && wrappedErrs[i-1] == argNum {
				continue
			}
			if e, ok := a[argNum].(error); ok {
				errs = append(errs, e)
			}
		}
		err = &wrapErrors{s, errs}
	}

	return err
}

// isSimpleFormat checks if the format string is just a single verb without flags
func isSimpleFormat(format string) bool {
	return format == "%s" || format == "%w" || format == "%v"
}

// parseErrorFormat parses a format string, replacing %w with %v
// and tracking the positions of wrapped errors
func parseErrorFormat(format string, wrappedErrs *[]int, reordered *bool) string {
	var sb strings.Builder
	argNum := 0

	for i := 0; i < len(format); {
		// Copy regular characters
		if format[i] != '%' {
			sb.WriteByte(format[i])
			i++
			continue
		}

		// We have a '%'
		if i+1 < len(format) && format[i+1] == '%' {
			// Escaped percent sign: %%
			sb.WriteString("%%")
			i += 2
			continue
		}

		// Beginning of a format verb
		var formatVerb strings.Builder
		formatVerb.WriteByte('%')
		i++

		// Check for reordering [n]
		if i < len(format) && format[i] == '[' {
			*reordered = true
			start := i
			for i < len(format) && format[i] != ']' {
				i++
			}
			if i < len(format) && format[i] == ']' {
				// Found the end of the index
				indexStr := format[start+1 : i]
				index := 0
				for j := 0; j < len(indexStr); j++ {
					if indexStr[j] >= '0' && indexStr[j] <= '9' {
						index = index*10 + int(indexStr[j]-'0')
					}
				}
				argNum = index
				formatVerb.WriteString(format[start : i+1])
				i++
			}
		}

		// Copy any flags, width, precision
		for i < len(format) {
			// Check for flags
			if strings.ContainsRune("+-0# ", rune(format[i])) {
				formatVerb.WriteByte(format[i])
				i++
				continue
			}

			// Check for width/precision
			if format[i] == '.' || (format[i] >= '0' && format[i] <= '9') || format[i] == '*' {
				formatVerb.WriteByte(format[i])
				i++
				continue
			}

			break
		}

		// Handle verb
		if i < len(format) {
			if format[i] == 'w' {
				// It's a %w verb - track this position and replace with %v
				*wrappedErrs = append(*wrappedErrs, argNum)
				formatVerb.WriteByte('v')
			} else {
				formatVerb.WriteByte(format[i])
			}
			i++
		}

		sb.WriteString(formatVerb.String())

		// Move to next argument
		argNum++
	}

	return sb.String()
}

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}

type wrapErrors struct {
	msg  string
	errs []error
}

func (e *wrapErrors) Error() string {
	return e.msg
}

func (e *wrapErrors) Unwrap() []error {
	return e.errs
}
