package bit_test

import (
	"testing"

	"github.com/rprtr258/sixty/encoding/bit"
	"github.com/rprtr258/sixty/internal"
	"github.com/stretchr/testify/assert"
)

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bit.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBits)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := bit.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBBits, got)
}
