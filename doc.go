/*
Package xprint provides high-performance string formatting functions that are API-compatible with fmt.

The package is designed to be a drop-in replacement for fmt.Sprintf with significant performance
improvements (40-47% faster in benchmarks). It achieves this by using specialized string building
techniques and optimized type-specific formatting paths.

Example usage:

	str := xprint.Printf("Hello, %s! You are %d years old.", "User", 25)

The package supports all standard formatting directives from the fmt package, including:
  - %v - default format
  - %+v - default format with struct field names
  - %#v - Go syntax format
  - %T - type
  - %t - boolean
  - %d - decimal integer
  - %b - binary integer
  - %c - character
  - %s - string
  - %p - pointer
  - And more

Performance improvements are most notable when:
  - Formatting strings with the %s verb
  - Working with a mix of string and numeric values
  - Handling complex data structures
*/
package xprint
