package stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/stream"
	"github.com/rprtr258/sixty/internal"
)

var bttvEmotesDBStream = []byte{
	0x01,
	0x05, 0x4b, 0x61, 0x70, 0x70, 0x61, // "Kappa"
	0x01,
	0x0b, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "KappaOwner1"
	0x01,
	0x0b, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "KappaOwner2"
	0x00,
	0x01,
	0x08, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, // "PogChamp"
	0x01,
	0x0e, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "PogChampOwner1"
	0x01,
	0x0e, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "PogChampOwner2"
	0x00,
	0x00,
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := stream.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDBStream, got)
}

func TestDecode_BttvDB(t *testing.T) {
	t.Parallel()

	type entry struct {
		Emote string
		Ids   []string
	}

	type BttvEmotesDB []entry

	var got BttvEmotesDB
	err := stream.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBStream, &got)
	assert.NoError(t, err)
	assert.Equal(t,
		BttvEmotesDB([]entry{{
			Emote: "Kappa",
			Ids:   []string{"KappaOwner1", "KappaOwner2"},
		}, {
			Emote: "PogChamp",
			Ids:   []string{"PogChampOwner1", "PogChampOwner2"},
		}}),
		got,
	)
}

func TestDecode_Coproduct(t *testing.T) {
	t.Parallel()

	type s any
	var got s
	err := stream.Decode(sixty.Coproduct{
		{"int", sixty.Int{}},
		{"string", sixty.String{}},
	}, []byte{0x00, 0b01010100}, &got)
	assert.NoError(t, err)
	assert.Equal(t, s(int64(42)), got)
}

func TestDecode_ArrayAndOptional(t *testing.T) {
	t.Parallel()

	type s struct {
		Arr  [2]uint64
		Opt  *string
		Opt2 *int64
	}

	var got s
	err := stream.Decode(sixty.Product{
		{"arr", sixty.Array{Schema: sixty.Uint{}, Count: 2}},
		{"opt", sixty.Optional{Schema: sixty.String{}}},
		{"opt2", sixty.Optional{Schema: sixty.Int{}}},
	}, []byte{0x01, 0x02, 0x01, 0x01, 0x78, 0x00}, &got)
	assert.NoError(t, err)
	opt := "x"
	assert.Equal(t, s{
		Arr:  [2]uint64{1, 2},
		Opt:  &opt,
		Opt2: nil,
	}, got)
}
