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
  std:string,
  []std:string,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_BttvDB(t *testing.T) {
	assert.Equal(t, `[](std:string*[]std:string)`, text.AlgebraicEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_VKInfo(t *testing.T) {
	assert.Equal(t, `(uint*std:string*std:string)+(uint*std:string*std:string)`, text.AlgebraicEncode(internal.VKInfoSchema))
}
