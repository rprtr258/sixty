package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bytes.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBytes)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := bytes.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBBytes, got)
}
