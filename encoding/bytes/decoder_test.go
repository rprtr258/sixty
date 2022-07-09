package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bytes "github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBDecode(t *testing.T) {
	got := bytes.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBytes)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}
