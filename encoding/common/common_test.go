package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/encoding/common"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	type entry struct {
		Emote string
		Ids   []string
	}

	type BttvEmotesDB []entry

	var got BttvEmotesDB
	err := common.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData, &got)
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

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	type entry struct {
		Emote string
		Ids   []string
	}

	got, err := common.Encode(internal.BttvEmotesDBSchema, []entry{{
		Emote: "Kappa",
		Ids:   []string{"KappaOwner1", "KappaOwner2"},
	}, {
		Emote: "PogChamp",
		Ids:   []string{"PogChampOwner1", "PogChampOwner2"},
	}})
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestDecode_Coproduct(t *testing.T) {
	t.Parallel()

	var got any
	err := common.Decode(
		sixty.Coproduct{
			{"int", sixty.Int{}},
			{"string", sixty.String{}},
		},
		sixty.CoproductValue{Variant: 0, SchemaValue: sixty.IntValue(42)},
		&got,
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), got)
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	type sample struct {
		Bit     bool
		Int     int64
		Uint    uint64
		Rune    rune
		Text    string
		Enum    uint8
		EnumInt int64
		Opt     *string
		List    []int64
		Arr     [2]int64
	}

	schema := sixty.Product{
		{"bit", sixty.Bit{}},
		{"int", sixty.Int{}},
		{"uint", sixty.Uint{}},
		{"rune", sixty.Rune{}},
		{"text", sixty.String{}},
		{"enum", sixty.Enum{Values: []string{"a", "b"}}},
		{"enum_int", sixty.EnumInt{Min: 10, Max: 12}},
		{"opt", sixty.Optional{Schema: sixty.String{}}},
		{"list", sixty.List{Schema: sixty.Int{}}},
		{"array", sixty.Array{Schema: sixty.Int{}, Count: 2}},
	}

	opt := "pog"
	for _, original := range []sample{{
		Bit:     true,
		Int:     -42,
		Uint:    42,
		Rune:    '☺',
		Text:    "κ산☺",
		Enum:    1,
		EnumInt: 11,
		Opt:     &opt,
		List:    []int64{1, -2, 3},
		Arr:     [2]int64{4, 5},
	}, {
		Bit:     false,
		Int:     0,
		Uint:    0,
		Rune:    'a',
		Text:    "",
		Enum:    0,
		EnumInt: 10,
		List:    []int64{},
		Arr:     [2]int64{0, 0},
	}} {
		value, err := common.Encode(schema, original)
		assert.NoError(t, err)
		data, err := bytes.Encode(schema, value)
		assert.NoError(t, err)
		value2, err := bytes.Decode(schema, data)
		assert.NoError(t, err)
		var got sample
		err = common.Decode(schema, value2, &got)
		assert.NoError(t, err)
		assert.Equal(t, original, got)
	}
}
