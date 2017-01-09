package strtool

import "unicode/utf8"

type Reader struct {
	b    []byte
	i, j int
	ch   rune // if RuneError then error or EOF
	sz   int
}

func NewTokenizer(b []byte) *Reader {
	var r = &Reader{
		b: b,
	}
	r.load()
	return r
}

func (r *Reader) load() {
	r.ch, r.sz = utf8.DecodeRune(r.b[r.i:])
}

func (r *Reader) read() {
	r.i += r.sz
}

func (r *Reader) ok() bool {
	return r.sz > 0
}

func (r *Reader) Init(b []byte) {
	r.b = b
	r.i = 0
	r.j = 0
	r.load()
}

// Pos returns the offset in bytes from the begining of the buffer.
func (r *Reader) Pos() int {
	return r.i
}

// Len returns the number of bytes of the unread portion of the buffer.
func (r *Reader) Len() int {
	return len(r.b) - r.i
}

func (r *Reader) Bytes() []byte {
	return r.b[r.j:r.i]
}

func (r *Reader) Discard() {
	r.j = r.i
}

func (r *Reader) Yield() string {
	bytes := r.Bytes()
	r.Discard()
	return string(bytes)
}

// Peek returns the currently loaded rune or utf8.RuneError.
func (r *Reader) Peek() rune {
	return r.ch
}

func (r *Reader) AcceptAny() bool {
	if r.ok() {
		r.read()
		r.load()
		return true
	}
	return false
}

func (r *Reader) AcceptRune(ch rune) bool {
	if r.ok() && r.ch == ch {
		r.read()
		r.load()
		return true
	}
	return false
}

func (r *Reader) AcceptSpace() bool {
	if r.ok() && isSpace(r.ch) {
		r.read()
		r.load()
		return true
	}
	return false
}

func (r *Reader) AcceptBetween(a, b rune) bool {
	if r.ok() && a <= r.ch && r.ch <= b {
		r.read()
		r.load()
		return true
	}
	return false
}

func (r *Reader) AcceptFunc(f func(ch rune) bool) bool {
	if r.ok() && f(r.ch) {
		r.read()
		r.load()
		return true
	}
	return false
}

type Snapshot struct {
	r Reader
}

func (r *Reader) Snapshot() Snapshot {
	return Snapshot{r: *r}
}

func (r *Reader) Backtrack(k Snapshot) {
	*r = k.r
}
