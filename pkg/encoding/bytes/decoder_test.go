package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root

	bytes "github.com/rprtr258/sixty/pkg/encoding/bytes"
	"github.com/rprtr258/sixty/tests"
)

func TestBttvDBDecode(t *testing.T) {
	got := bytes.Decode(tests.BttvEmotesDBSchema, tests.BttvEmotesDBBytes)
	assert.Equal(t, tests.BttvEmotesDBData, got)
}
