package sixty

import (
	"io"
	"log"
	"unicode/utf8"
)

const defaultBufferSize = 1 << 10

type Reader struct {
	reader io.Reader
	buffer []byte
	start  int
	end    int
	pos    uint
	Error  error
}

// TODO: test NewReader(bytes.NewStringBuffer(""))
func NewReaderBufferSize(r io.Reader, bufferSize int) Reader {
	return Reader{
		reader: r,
		buffer: make([]byte, bufferSize),
		Error:  nil,
		start:  0,
		end:    0,
	}
}

func NewReader(r io.Reader) Reader {
	return NewReaderBufferSize(r, defaultBufferSize)
}

func (r *Reader) loadMore() bool {
	for {
		n, err := r.reader.Read(r.buffer)
		log.Println(n, err)
		if n != 0 {
			r.start = 0
			r.end = n
			if err != nil {
				r.Error = err
			}
			return true
		}
		if err != nil {
			r.Error = err
			return false
		}
	}
}

func (r *Reader) Pos() uint {
	return r.pos
}

func (r *Reader) ReadByte() (b byte) {
	if r.Error != nil {
		return
	}
	if r.start == r.end {
		r.loadMore()
		if r.Error != nil {
			return
		}
	}
	b = r.buffer[r.start]
	r.start++
	r.pos++
	return
}

func (r *Reader) ReadRune() (u rune) {
	if r.Error != nil {
		return
	}
	if r.start == r.end {
		if !r.loadMore() {
			return
		}
	}
	// TODO: rune is not in buffer, or just implement decode rune myself using ReadByte
	u, n := utf8.DecodeRune(r.buffer[r.start:])
	// TODO: error reading rune
	r.start += n
	r.pos += uint(n)
	return
}
