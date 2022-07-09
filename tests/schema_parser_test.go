package sixty_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root
	sixty "github.com/rprtr258/sixty/pkg"
	sixtyBytes "github.com/rprtr258/sixty/pkg/byte"
	"github.com/rprtr258/sixty/pkg/json"
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
	buffer := bytes.Buffer{}
	(&json.JsonEncoder{Buffer: &buffer}).Write(bttvEmotesDBSchema, bttvEmotesDBData)
	assert.Equal(t, bttvEmotesDBJson, buffer.String())
}

func TestBttvDBByteEncoder(t *testing.T) {
	buffer := bytes.Buffer{}
	(&sixtyBytes.ByteEncoder{Buffer: &buffer}).Write(bttvEmotesDBSchema, bttvEmotesDBData)
	assert.Equal(t, bttvEmotesDBBytes, buffer.Bytes())
}

func TestBttvDBByteDecoder(t *testing.T) {
	buffer := bytes.NewBuffer(bttvEmotesDBBytes)
	got := (&sixtyBytes.ByteDecoder{Buffer: buffer}).Read(bttvEmotesDBSchema)
	assert.Equal(t, bttvEmotesDBData, got)
}

func TestBttvDBSize(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SchemaSize(bttvEmotesDBSchema))
}

func TestBttvDB(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](std:string,[]std:string)`)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDBSchema, schema)
}

func TestBttvDBPrettyPrinter(t *testing.T) {
	assert.Equal(t, `[](
  std:string,
  []std:string,
)`, sixty.PrettyPrintSchema(bttvEmotesDBSchema))
}

func TestEnumInt(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](0..255)`)
	assert.NoError(t, err)
	assert.Equal(t, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}}, schema)

}
