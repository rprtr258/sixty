package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func assertSchema(t *testing.T, schemaString string, expected sixty.Schema) {
	got, err := text.ParseSchema(schemaString)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestDecode_BttvDB(t *testing.T) {
	assertSchema(t, internal.BttvEmotesDBSchemaString, internal.BttvEmotesDBSchema)
}

func TestDecode_VKInfo(t *testing.T) {
	assertSchema(t, internal.VKInfoSchemaString, internal.VKInfoSchema)
}

func TestDecode_EnumInt(t *testing.T) {
	assertSchema(t, `[][0..255]`, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}})
}
