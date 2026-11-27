package json_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/json"
	"github.com/rprtr258/sixty/internal"
)

//go:embed bttv_emotes_db.json
var bttvEmotesDBJson string

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := json.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDBJson, got)
}
