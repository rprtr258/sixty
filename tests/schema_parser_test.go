package sixty_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	sixty "github.com/rprtr258/sixty/pkg"
	bytesenc "github.com/rprtr258/sixty/pkg/byte"
)

var (
	bttvEmotesDB = sixty.List{Type: sixty.Product{
		sixty.String, // TODO: name as "emote"
		sixty.List{
			// TODO: change to const(?) length id
			Type: sixty.String,
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
		0x2,                                                                   // array of 2 strings
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "KappaOwner1"
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "KappaOwner2"
		0x8, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, // "PogChamp"
		0x2,
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "PogChampOwner1"
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "PogChampOwner2"
	}
)

// TODO: use to encode string and other primitive types
func stringToString(s string) sixty.TypeValue {
	var l sixty.ListValue
	for _, c := range s {
		l = append(l, sixty.RuneValue(c))
	}
	return sixty.AliasValue{
		Name:  sixty.String.Name,
		Value: l,
	}
}

func TestBttvDBJsonEncoder(t *testing.T) {
	je := &sixty.JsonEncoder{}
	sixty.VisitValue(bttvEmotesDB, bttvEmotesDBData, je)
	assert.Equal(t, `[{"0":"Kappa","1":["KappaOwner1","KappaOwner2"]},{"0":"PogChamp","1":["PogChampOwner1","PogChampOwner2"]}]`, je.String())
}

func TestBttvDBByteEncoder(t *testing.T) {
	je := bytes.Buffer{}
	sixty.VisitValue(bttvEmotesDB, bttvEmotesDBData, &bytesenc.ByteEncoder{Builder: &je})
	assert.Equal(t, bttvEmotesDBBytes, je.Bytes())
}

func TestBttvDBByteDecoder(t *testing.T) {
	je := &bytesenc.ByteDecoder{Buf: bytes.NewBuffer(bttvEmotesDBBytes)}
	got := sixty.VisitValue2(bttvEmotesDB, je)
	assert.Equal(t, bttvEmotesDBData, got)
}

func TestBttvDBSize(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SizeVisitor(bttvEmotesDB))
}

func TestBttvDB(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](std:string,[]std:string)`)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDB, schema)
}

func TestBttvDBPrettyPrinter(t *testing.T) {
	assert.Equal(t, `std:string = []rune

[](
  std:string,
  []std:string,
)`, sixty.PrettyPrintSchema(bttvEmotesDB))
}

func TestEnumInt(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](0..255)`)
	assert.NoError(t, err)
	assert.Equal(t, sixty.List{Type: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}}, schema)

}
