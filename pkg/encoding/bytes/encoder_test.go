package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root

	"github.com/rprtr258/sixty/pkg/encoding/bytes"
	"github.com/rprtr258/sixty/tests"
)

func TestBttvDBEncode(t *testing.T) {
	got := bytes.Encode(tests.BttvEmotesDBSchema, tests.BttvEmotesDBData)
	assert.Equal(t, tests.BttvEmotesDBBytes, got)
}
