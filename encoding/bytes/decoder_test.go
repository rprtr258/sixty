package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bytes "github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBDecode(t *testing.T) {
	got, err := bytes.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBytes)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}
