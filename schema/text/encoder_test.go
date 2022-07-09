package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func TestBttvDBEncode(t *testing.T) {
	assert.Equal(t, `[](std:string,[]std:string)`, text.Encode(internal.BttvEmotesDBSchema))
}

func TestBttvDBPrettyEncode(t *testing.T) {
	assert.Equal(t, `[](
  std:string,
  []std:string,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}
