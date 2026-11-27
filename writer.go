package sixty

import (
	"io"
	"unicode/utf8"

	"github.com/rprtr258/sixty/internal/zigzag"
)

const defaultWriteBufferSize = 1 << 10

type Writer struct {
	writer io.Writer
	// buffer - bytes written but not yet flushed to the underlying writer
	buffer     []byte
	runeBuffer []byte
	Error      error
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		writer:     w,
		buffer:     make([]byte, 0, defaultWriteBufferSize),
		runeBuffer: make([]byte, 4),
	}
}

// ensure makes room for n more bytes in the buffer, flushing to the underlying
// writer when it is full, and reports whether writing can proceed.
func (w *Writer) ensure(n int) bool {
	if w.Error != nil {
		return false
	}
	if len(w.buffer)+n <= cap(w.buffer) {
		return true
	}
	w.flush()
	return w.Error == nil
}

// flush writes buffered data to the underlying writer.
func (w *Writer) flush() {
	if w.Error != nil || len(w.buffer) == 0 {
		return
	}
	if _, err := w.writer.Write(w.buffer); err != nil {
		w.Error = err
	}
	w.buffer = w.buffer[:0]
}

// Flush writes buffered data to the underlying writer.
func (w *Writer) Flush() {
	w.flush()
}

// write appends b to the buffer, or writes it directly if it does not fit.
func (w *Writer) write(b []byte) {
	if !w.ensure(len(b)) {
		return
	}
	if len(b) > cap(w.buffer) {
		if _, err := w.writer.Write(b); err != nil {
			w.Error = err
		}
		return
	}
	w.buffer = append(w.buffer, b...)
}

func (w *Writer) WriteByte(b byte) {
	w.write([]byte{b})
}

func (w *Writer) WriteRune(r rune) {
	n := utf8.EncodeRune(w.runeBuffer, r)
	w.write(w.runeBuffer[:n])
}

func (w *Writer) WriteString(s string) {
	if !w.ensure(len(s)) {
		return
	}
	if len(s) > cap(w.buffer) {
		if _, err := io.WriteString(w.writer, s); err != nil {
			w.Error = err
		}
		return
	}
	w.buffer = append(w.buffer, s...)
}

func (w *Writer) WriteUint64(x uint64) {
	for x > 0x7F {
		w.WriteByte(byte(x&0x7F) | 0x80)
		x >>= 7
	}
	w.WriteByte(byte(x))
}

func (w *Writer) WriteInt64(x int64) {
	w.WriteUint64(zigzag.Encode(x))
}
