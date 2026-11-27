package sixty

import (
	"io"
	"unicode/utf8"
)

type Writer struct {
	writer io.Writer
	// TODO: buffer
	runeBuffer []byte
	Error      error
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		writer:     w,
		runeBuffer: make([]byte, 4),
		Error:      nil,
	}
}

func (w *Writer) WriteByte(b byte) {
	if w.Error != nil {
		return
	}
	if _, err := w.writer.Write([]byte{b}); err != nil {
		w.Error = err
	}
	return
}

func (w *Writer) WriteRune(r rune) {
	if w.Error != nil {
		return
	}
	n := utf8.EncodeRune(w.runeBuffer, r)
	if _, err := w.writer.Write(w.runeBuffer[:n]); err != nil {
		w.Error = err
	}
}

func (w *Writer) WriteString(s string) {
	if w.Error != nil {
		return
	}
	if _, err := w.writer.Write([]byte(s)); err != nil {
		w.Error = err
	}
}

func (w *Writer) WriteUint64(x uint64) {
	for x > 0x7F {
		w.WriteByte(byte(x&0x7F) | 0x80)
		x >>= 7
	}
	w.WriteByte(byte(x))
}

func (w *Writer) WriteInt64(x int64) {
	w.WriteUint64(ToZigZag(x))
}
