package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func TestEncode_BttvDB(t *testing.T) {
	assert.Equal(t, internal.BttvEmotesDBSchemaString, text.Encode(internal.BttvEmotesDBSchema))
}

func TestEncode_VKInfo(t *testing.T) {
	assert.Equal(t, internal.VKInfoSchemaString, text.Encode(internal.VKInfoSchema))
}

func TestPrettyEncode_BttvDB(t *testing.T) {
	assert.Equal(t, `[](
  emote = :string,
  ids = []:string,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_BttvDB(t *testing.T) {
	assert.Equal(t, `[](:string * []:string)`, text.AlgebraicEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_VKInfo(t *testing.T) {
	assert.Equal(t, `(uint * :string * :string) + (uint * :string * :string)`, text.AlgebraicEncode(internal.VKInfoSchema))
}
