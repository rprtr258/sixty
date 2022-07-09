package bit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/bit"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	got, err := bit.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBBits, got)
}
