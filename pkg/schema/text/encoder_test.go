package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root

	"github.com/rprtr258/sixty/pkg/schema/text"
	"github.com/rprtr258/sixty/tests"
)

func TestBttvDBEncode(t *testing.T) {
	assert.Equal(t, `[](
  std:string,
  []std:string,
)`, text.Encoder(tests.BttvEmotesDBSchema))
}
