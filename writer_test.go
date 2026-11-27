package sixty_test

import (
	"bytes"
	"errors"
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
		require.Equal(t, test.encoding, buffer.Bytes())
	}
}

func TestWriteErrorByte(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteByte(1)
	require.Error(t, writer.Error)
}

func TestWriteErrorRune(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteRune('ж')
	require.Error(t, writer.Error)
}

func TestWriteErrorString(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteString("hello")
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
	require.Error(t, writer.Error)
}

func TestUint64RoundTrip(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	writer := sixty.NewWriter(&buffer)
	writer.WriteUint64(300)
	reader := sixty.NewReader(&buffer)
	require.Equal(t, uint64(300), reader.ReadUint64())
	require.NoError(t, reader.Error)
}

func TestInt64RoundTrip(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	writer := sixty.NewWriter(&buffer)
	writer.WriteInt64(-42)
	reader := sixty.NewReader(&buffer)
	require.Equal(t, int64(-42), reader.ReadInt64())
	require.NoError(t, reader.Error)
}
