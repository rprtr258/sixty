package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/common/json"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := json.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBJson, got)
}
