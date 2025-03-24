# xprint validation

This package contains the validation and benchmarking tools for the xprint package. It is not part of the exported API and is only used for development and testing.

## Structure

- `main.go` - Entry point for validation CLI
- `simple.go` - API compatibility tests comparing to fmt.Sprintf
- `newbench.go` - Performance benchmarking suite
- `runtest.go` - Additional test utilities
- `internal/largeints` - Test data utilities

## Usage

First, build the validation binary:

```bash
go build -o validation.bin .
```

Run the different validation tools:

```bash
# Run compatibility tests
./validation.bin simple

# Run performance benchmarks
./validation.bin newbench

# Run quick tests
./validation.bin quick

# Run mixed type tests
./validation.bin mixedtype
```

## CI Integration

This validation suite is automatically run in the CI pipeline to ensure compatibility and performance standards are maintained. 