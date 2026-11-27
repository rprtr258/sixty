package sixty

import (
	"fmt"
	"io"
)

const defaultBufferSize = 1 << 10

// TooLongRuneError - reported by ReadRune on a UTF-8 lead byte claiming more than 4 bytes.
type TooLongRuneError struct {
	Pos uint
}

func (e TooLongRuneError) Error() string {
	return fmt.Sprintf("TOO LONG RUNE at %d", e.Pos)
}

type Reader struct {
	// reader - underlying raw reader
	reader io.Reader
	// buffer - buffer for data
	buffer []byte
	// start - first position of data in buffer
	start int
	// end - last position of data in buffer plus one
	end int
	// pos - count of bytes scanned
	pos uint
	// Error - first error occurred
	Error error
}

func NewReaderBufferSize(reader io.Reader, bufferSize uint) Reader {
	return Reader{
		reader: reader,
		buffer: make([]byte, bufferSize),
	}
}

func NewReader(r io.Reader) Reader {
	return NewReaderBufferSize(r, defaultBufferSize)
}

// loadMore - read data from inner reader and update buffer, returns true if some data was read.
func (r *Reader) loadMore() bool {
	if r.start != r.end {
		return true
	}
	if r.Error != nil {
		return false
	}
	for {
		n, err := r.reader.Read(r.buffer)
		if err != io.EOF {
			r.Error = err
		}
		if n > 0 {
			r.start = 0
			r.end = n
			return true
		}
		if r.Error != nil || err == io.EOF {
			return false
		}
	}
}

func (r *Reader) Pos() uint {
	return r.pos
}

func (r *Reader) ReadByte() (b byte) {
	if r.Error != nil || !r.loadMore() {
		return
	}
	b = r.buffer[r.start]
	r.start++
	r.pos++
	return
}

func (r *Reader) ReadRune() (u rune) {
	const maskx = 0x3F
	b0 := r.ReadByte()
	switch leadingOnes(b0) {
	case 0:
		return rune(b0)
	case 1:
		b1 := r.ReadByte()
		return rune(b0&0x3F)<<6 | rune(b1&maskx)
	case 2:
		b1 := r.ReadByte()
		b2 := r.ReadByte()
		return rune(b0&0x1F)<<12 | rune(b1&maskx)<<6 | rune(b2&maskx)
	case 3:
		b1 := r.ReadByte()
		b2 := r.ReadByte()
		b3 := r.ReadByte()
		return rune(b0&0x0F)<<18 | rune(b1&maskx)<<12 | rune(b2&maskx)<<6 | rune(b3&maskx)
	default:
		if r.Error == nil {
			r.Error = TooLongRuneError{Pos: r.pos}
		}
		return
	}
}

// ReadLengthString - read length encoded string.
func (r *Reader) ReadLengthString() string {
	resLen := r.ReadUint64()
	runes := make([]rune, 0, resLen)
	for i := uint64(0); i < resLen; i++ {
		runes = append(runes, r.ReadRune())
	}
	return string(runes)
}

func (r *Reader) ReadUint64() (u uint64) {
	i := 0
	for {
		b := r.ReadByte()
		u |= uint64(b&0x7F) << (7 * i)
		if b&0x80 == 0 {
			return
		}
		i++
	}
}

func (r *Reader) ReadInt64() int64 {
	return FromZigZag(r.ReadUint64())
}

func leadingOnes(x byte) (res uint) {
	for ; x&0x80 != 0; x <<= 1 {
		res++
	}
	return
}
