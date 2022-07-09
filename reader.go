package sixty

import (
	"io"
	"log"
	"unicode/utf8"
)

const readerBufferSize = 1 << 10

type Reader struct {
	reader io.Reader
	buffer []byte
	start  int
	end    int
	Error  error
}

// TODO: test NewReader(bytes.NewStringBuffer(""))
func NewReader(r io.Reader) Reader {
	return Reader{
		reader: r,
		buffer: make([]byte, readerBufferSize),
		Error:  nil,
		start:  0,
		end:    0,
	}
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
	return
}
