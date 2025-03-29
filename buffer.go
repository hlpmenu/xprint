package xprint

import (
	"strings"
	"unicode/utf8"
)

// buffer is a simple []byte buffer for building strings.
type buffer []byte

func (b *buffer) Len() int {
	return len([]byte(*b))
}
func (b *buffer) LenMB() int {
	return BtoMB(b.Len())
}

func (b *buffer) write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}
func (b *buffer) writeStringToUpper(s string) {
	*b = append(*b, strings.ToUpper(s)...)
}

func (b *buffer) writeByte(c byte) {
	*b = append(*b, c)
}

func (b *buffer) writeRune(r rune) {
	*b = utf8.AppendRune(*b, r)
}
func BtoMB(b int) int {
	return b / 1024 / 1024
}

func (b *buffer) writeNilArg(verb rune) {
	*b = append(*b, percentBangString...)
	*b = append(*b, byte(verb))
	*b = append(*b, '(')
	*b = append(*b, nilAngleString...)
	*b = append(*b, ')')

}
