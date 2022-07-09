package sixty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root
	sixty "github.com/rprtr258/sixty/pkg"
	sixtyBytes "github.com/rprtr258/sixty/pkg/encoding/byte"
	sixtyJson "github.com/rprtr258/sixty/pkg/encoding/json"
	textSchema "github.com/rprtr258/sixty/pkg/schema/text"
)

var (
	bttvEmotesDBSchema = sixty.List{Schema: sixty.Product{
		sixty.String{}, // TODO: name as "emote"
		sixty.List{
			// TODO: change to const(?) length id
			Schema: sixty.String{},
		}, // TODO: name as "ids"
	}}
	bttvEmotesDBData = sixty.ListValue{sixty.ProductValue{
		stringToString("Kappa"),
		sixty.ListValue{
			stringToString("KappaOwner1"),
			stringToString("KappaOwner2"),
		},
	}, sixty.ProductValue{
		stringToString("PogChamp"),
		sixty.ListValue{
			stringToString("PogChampOwner1"),
			stringToString("PogChampOwner2"),
		},
	}}
	bttvEmotesDBBytes = []byte{
		0x2,
		0x5, 0x4b, 0x61, 0x70, 0x70, 0x61, // "Kappa"
		0x2,
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "KappaOwner1"
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "KappaOwner2"
		0x8, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, // "PogChamp"
		0x2,
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "PogChampOwner1"
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "PogChampOwner2"
	}
	bttvEmotesDBJson = `[{"0":"Kappa","1":["KappaOwner1","KappaOwner2"]},{"0":"PogChamp","1":["PogChampOwner1","PogChampOwner2"]}]`
)

// TODO: use to encode string and other primitive types
func stringToString(s string) sixty.SchemaValue {
	return sixty.StringValue(s)
}

func TestBttvDBJsonEncoder(t *testing.T) {
	got := sixtyJson.Encode(bttvEmotesDBSchema, bttvEmotesDBData)
	assert.Equal(t, bttvEmotesDBJson, got)
}

func TestBttvDBByteEncoder(t *testing.T) {
	got := sixtyBytes.Encode(bttvEmotesDBSchema, bttvEmotesDBData)
	assert.Equal(t, bttvEmotesDBBytes, got)
}

func TestBttvDBByteDecoder(t *testing.T) {
	got := sixtyBytes.Decode(bttvEmotesDBSchema, bttvEmotesDBBytes)
	assert.Equal(t, bttvEmotesDBData, got)
}

func TestBttvDBSize(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SchemaSize(bttvEmotesDBSchema))
}

func TestBttvDB(t *testing.T) {
	schema, err := textSchema.ParseSchema(`[](std:string,[]std:string)`)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDBSchema, schema)
}

func TestBttvDBPrettyPrinter(t *testing.T) {
	assert.Equal(t, `[](
  std:string,
  []std:string,
)`, textSchema.Encoder(bttvEmotesDBSchema))
}

func TestEnumInt(t *testing.T) {
	schema, err := textSchema.ParseSchema(`[](0..255)`)
	assert.NoError(t, err)
	assert.Equal(t, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}}, schema)

}
