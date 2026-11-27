package sixty_test

import (
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

func TestBufferedErrorPropagation(t *testing.T) {
	t.Parallel()

	writer := sixty.NewWriter(errWriter{errDiskFull})
	writer.WriteString(strings.Repeat("a", 2048)) // forces a flush to the failing writer
	require.ErrorIs(t, writer.Error, errDiskFull)
}
