package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root

	"github.com/rprtr258/sixty/pkg/encoding/json"
	"github.com/rprtr258/sixty/tests"
)

func TestBttvDBEncode(t *testing.T) {
	got := json.Encode(tests.BttvEmotesDBSchema, tests.BttvEmotesDBData)
	assert.Equal(t, tests.BttvEmotesDBJson, got)
}
