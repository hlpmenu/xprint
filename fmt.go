package xprint

// fmtFlags contains the core formatting flags
type fmtFlags struct {
	widPresent, precPresent         bool
	minus, plus, sharp, space, zero bool
	plusV, sharpV                   bool
	wid, prec                       int
	uintbase                        int
	toupper                         bool
}

// fmt holds the formatting state
type fmt struct {
	buf *buffer
	fmtFlags
	// intbuf is large enough to store %b of an int64 with a sign and
	// avoids padding at the end of the struct on 32 bit architectures.
}

func (f *fmt) init(b *buffer) {
	f.buf = b
	f.clearflags()
}

func (f *fmt) clearflags() {
	f.widPresent = false
	f.precPresent = false
	f.minus = false
	f.plus = false
	f.sharp = false
	f.space = false
	f.zero = false
	f.plusV = false
	f.sharpV = false
	f.wid = 0
	f.prec = 0
}

// visited tracks pointers already seen during recursive value formatting
type visited struct {
	ptrs map[uintptr]bool
}

func (v *visited) init() {
	v.ptrs = make(map[uintptr]bool)
}

func (v *visited) visit(p uintptr) bool {
	if v.ptrs[p] {
		return true
	}
	v.ptrs[p] = true
	return false
}
