package bit_test

import (
	"testing"

	"github.com/rprtr258/sixty/encoding/common/bit"
	"github.com/rprtr258/sixty/internal"
	"github.com/stretchr/testify/assert"
)

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bit.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBits)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}
