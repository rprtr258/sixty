package sixty

import (
	"io"
	"unicode/utf8"
)

type Writer struct {
	writer io.Writer
	// TODO: buffer
	runeBuffer []byte
	err        error
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		writer:     w,
		runeBuffer: make([]byte, 4),
		err:        nil,
	}
}

func (w *Writer) WriteByte(b byte) {
	if w.err != nil {
		return
	}
	if _, err := w.writer.Write([]byte{b}); err != nil {
		w.err = err
	}
}

func (w *Writer) WriteString(s string) {
	if w.err != nil {
		return
	}
	if _, err := w.writer.Write([]byte(s)); err != nil {
		w.err = err
	}
}

func (w *Writer) WriteRune(r rune) {
	if w.err != nil {
		return
	}
	n := utf8.EncodeRune(w.runeBuffer, r)
	if _, err := w.writer.Write(w.runeBuffer[:n]); err != nil {
		w.err = err
	}
}
