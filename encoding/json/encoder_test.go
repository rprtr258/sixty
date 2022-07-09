package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/json"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	got := json.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.Equal(t, internal.BttvEmotesDBJson, got)
}
