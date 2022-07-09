package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func TestBttvDBDecode(t *testing.T) {
	schema, err := text.ParseSchema(`[](std:string,[]std:string)`)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBSchema, schema)
}

func TestEnumIntDecode(t *testing.T) {
	schema, err := text.ParseSchema(`[](0..255)`)
	assert.NoError(t, err)
	assert.Equal(t, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}}, schema)

}
