package sixty_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
)

// errDiskFull is the write error returned by errWriter.
var errDiskFull = errors.New("disk full")

// errWriter fails every write with the given error.
type errWriter struct{ err error }

func (w errWriter) Write([]byte) (int, error) { return 0, w.err }

func TestWritePrimitives(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	writer := sixty.NewWriter(&buffer)
	writer.WriteByte(1)
	writer.WriteRune('ж')
	writer.WriteString("hello")
	require.NoError(t, writer.Error)
	writer.Flush()
	require.Equal(t, []byte{1, 0xD0, 0xB6, 'h', 'e', 'l', 'l', 'o'}, buffer.Bytes())
}

func TestWriteUint64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    uint64
		encoding []byte
	}{
		// LEB128-style: 7 bits per byte, high bit marks continuation
		{0, []byte{0}},
		{1, []byte{1}},
		{0x7F, []byte{0x7F}},
		{0x80, []byte{0x80, 0x01}},
		{300, []byte{0xAC, 0x02}},
	}
	for _, test := range testCases {
		var buffer bytes.Buffer
		writer := sixty.NewWriter(&buffer)
		writer.WriteUint64(test.value)
		require.NoError(t, writer.Error)
		writer.Flush()
		require.Equal(t, test.encoding, buffer.Bytes())
	}
}

func TestWriteInt64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    int64
		encoding []byte
	}{
		{0, []byte{0}},
		{-1, []byte{1}},
		{1, []byte{2}},
		{-42, []byte{0x53}}, // zigzag(-42) = 83
		{42, []byte{0x54}},
	}
	for _, test := range testCases {
		var buffer bytes.Buffer
		writer := sixty.NewWriter(&buffer)
		writer.WriteInt64(test.value)
		require.NoError(t, writer.Error)
		writer.Flush()
		require.Equal(t, test.encoding, buffer.Bytes())
	}
}

func TestWriteErrorByte(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteByte(1)
	writer.Flush()
	require.Error(t, writer.Error)
}

func TestWriteErrorRune(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteRune('ж')
	writer.Flush()
	require.Error(t, writer.Error)
}

func TestWriteErrorString(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteString("hello")
	writer.Flush()
	require.Error(t, writer.Error)
}

func TestWritesSkippedAfterError(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteByte(1)
	writer.WriteRune('ж')
	writer.WriteString("hello")
	writer.WriteUint64(300)
	writer.WriteInt64(-1)
	writer.Flush()
	require.Error(t, writer.Error)
}

func TestUint64RoundTrip(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	writer := sixty.NewWriter(&buffer)
	writer.WriteUint64(300)
	writer.Flush()
	reader := sixty.NewReader(&buffer)
	require.Equal(t, uint64(300), reader.ReadUint64())
	require.NoError(t, reader.Error)
}

func TestInt64RoundTrip(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	writer := sixty.NewWriter(&buffer)
	writer.WriteInt64(-42)
	writer.Flush()
	reader := sixty.NewReader(&buffer)
	require.Equal(t, int64(-42), reader.ReadInt64())
	require.NoError(t, reader.Error)
}

// spyWriter records every Write call separately.
type spyWriter struct {
	writes [][]byte
}

func (w *spyWriter) Write(b []byte) (int, error) {
	copyBytes := make([]byte, len(b))
	copy(copyBytes, b)
	w.writes = append(w.writes, copyBytes)
	return len(b), nil
}

func TestWritesBufferedUntilFlush(t *testing.T) {
	t.Parallel()

	spy := &spyWriter{}
	writer := sixty.NewWriter(spy)
	writer.WriteString("hello")
	require.Empty(t, spy.writes)
	writer.Flush()
	require.Equal(t, [][]byte{[]byte("hello")}, spy.writes)
}

func TestBufferOverflowFlushesAutomatically(t *testing.T) {
	t.Parallel()

	spy := &spyWriter{}
	writer := sixty.NewWriter(spy)
	// default buffer is 1<<10 bytes; writing more forces an intermediate flush
	big := strings.Repeat("a", 2048)
	writer.WriteString(big)
	writer.WriteByte('b')
	writer.Flush()
	require.Equal(t, []byte(big), spy.writes[0])
	require.Equal(t, [][]byte{[]byte(big), []byte("b")}, spy.writes)
}

func TestLargeStringWrittenDirectly(t *testing.T) {
	t.Parallel()

	spy := &spyWriter{}
	writer := sixty.NewWriter(spy)
	writer.WriteByte('x')
	big := strings.Repeat("b", 2048)
	writer.WriteString(big)
	require.Equal(t, [][]byte{[]byte("x"), []byte(big)}, spy.writes)
	require.NoError(t, writer.Error)
}

func TestBufferedErrorPropagation(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteString(strings.Repeat("a", 2048)) // forces a flush to the failing writer
	require.ErrorIs(t, writer.Error, errDiskFull)
}
