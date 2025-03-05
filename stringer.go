package xprint

// Stringer is implemented by any value that has a String method.
type Stringer interface {
	String() string
}

// GoStringer is implemented by any value that has a GoString method.
type GoStringer interface {
	GoString() string
}
