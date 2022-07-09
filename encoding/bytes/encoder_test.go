package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	got := bytes.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.Equal(t, internal.BttvEmotesDBBytes, got)
}
