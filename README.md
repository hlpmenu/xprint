# xprint

A high-performance string formatting package for Go that provides significant performance improvements over the standard `fmt` package.

## Features

- **High Performance**: 40-47% faster than `fmt.Sprintf` in benchmark tests
- **API Compatible**: Drop-in replacement for `fmt.Sprintf` with the same formatting verbs
- **Specialized String Handling**: Optimized for common string formatting scenarios

## Installation

```bash
go get gopkg.hlmpn.dev/pkg/xprint
```

## Usage

```go
import (
    "gopkg.hlmpn.dev/pkg/xprint"
)

func main() {
    // Use like fmt.Sprintf
    str := xprint.Printf("Hello, %s! You are %d years old.", "User", 25)
    
    // Works with all standard fmt verbs
    complexStr := xprint.Printf("%v %#v %T", myStruct, myMap, myInterface)
}
```

## Performance

Benchmarks show that `xprint.Printf` is approximately 40-47% faster than `fmt.Sprintf` across a variety of use cases:

| Function | Average Time | ns/byte |
|----------|--------------|---------|
| xprint.Printf | 7.8ms | 0.07 |
| fmt.Sprintf | 13.4ms | 0.12 |

The performance advantage is most notable when:
- Formatting strings with the `%s` verb
- Working with a mix of string and numeric values
- Handling complex data structures

## Implementation Details

The performance improvement comes from:
- Using `strings.Builder` for efficient string concatenation
- Specialized fast paths for common formatting patterns
- Optimized number-to-string conversions

## Project Structure

- `/` - Main package files (exported API)
- `/validation` - Test and benchmark suite (not part of the exported API)
  - `/validation/internal` - Internal utilities for testing and benchmarking

## Validation

The validation suite in `/validation` provides comprehensive compatibility testing with `fmt.Sprintf`:

```bash
cd validation
go build -o validation.bin .
./validation.bin simple  # Run compatibility tests
./validation.bin newbench  # Run benchmarks
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 