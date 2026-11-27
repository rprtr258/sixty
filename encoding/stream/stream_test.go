package stream_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	require.Equal(t, bttvEmotesDBStream, got)
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
	require.NoError(t, err)
	require.Equal(t,
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

func TestNonAsciiStringRoundTrip(t *testing.T) {
	t.Parallel()

	// length prefix counts runes: "ёжик" is 5 runes but 9 bytes
	data, err := stream.Encode(sixty.String{}, sixty.StringValue("ёжик"))
	require.NoError(t, err)

	var got string
	require.NoError(t, stream.Decode(sixty.String{}, data, &got))
	require.Equal(t, "ёжик", got)
}

func TestDecode_Coproduct(t *testing.T) {
	t.Parallel()

	type s any
	var got s
	err := stream.Decode(sixty.Coproduct{
		{"int", sixty.Int{}},
		{"string", sixty.String{}},
	}, []byte{0x00, 0b01010100}, &got)
	require.NoError(t, err)
	require.Equal(t, s(int64(42)), got)
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
	require.NoError(t, err)
	opt := "x"
	require.Equal(t, s{
		Arr:  [2]uint64{1, 2},
		Opt:  &opt,
		Opt2: nil,
	}, got)
}

type inner struct {
	X, Y int64
}

type allTypes struct {
	Bit     bool
	Int     int64
	Uint    uint64
	Rune    rune
	Enum    uint8
	EnumInt int64
	Opt     *int64
	Opt2    int64 // non-pointer optional target
	Arr     [2]bool
	Slice   []bool // array into slice target
	Prod    inner
	Cop     any
	List    []int64
	Str     string
}

var allTypesSchema = sixty.Product{
	{"bit", sixty.Bit{}},
	{"i", sixty.Int{}},
	{"uint", sixty.Uint{}},
	{"rune", sixty.Rune{}},
	{"enum", sixty.Enum{Values: []string{"a", "b"}}},
	{"enum_int", sixty.EnumInt{Min: -5, Max: 5}},
	{"opt", sixty.Optional{Schema: sixty.Int{}}},
	{"opt2", sixty.Optional{Schema: sixty.Int{}}},
	{"arr", sixty.Array{Schema: sixty.Bit{}, Count: 2}},
	{"slice", sixty.Array{Schema: sixty.Bit{}, Count: 2}},
	{"prod", sixty.Product{{"x", sixty.Int{}}, {"y", sixty.Int{}}}},
	{"cop", sixty.Coproduct{{"i", sixty.Int{}}, {"string", sixty.String{}}}},
	{"list", sixty.List{Schema: sixty.Int{}}},
	{"str", sixty.String{}},
}

func TestAllTypesRoundTrip(t *testing.T) {
	t.Parallel()

	schema := allTypesSchema
	data := allTypes{
		Bit:     true,
		Int:     -42,
		Uint:    300,
		Rune:    'ж',
		Enum:    1,
		EnumInt: -3,
		Opt:     nil,
		Opt2:    9,
		Arr:     [2]bool{true, false},
		Slice:   []bool{true, false},
		Prod:    inner{X: 1, Y: -2},
		Cop:     "hi",
		List:    []int64{1, -2},
		Str:     "hello",
	}

	// stream.Encode takes SchemaValues, so data is spelled as a ProductValue here
	enc, err := stream.Encode(schema, sixty.ProductValue{
		sixty.BitValue(true),
		sixty.IntValue(-42),
		sixty.UintValue(300),
		sixty.RuneValue('ж'),
		sixty.EnumValue(1),
		sixty.EnumIntValue(-3),
		sixty.OptionalValue{IsValid: false},
		sixty.OptionalValue{SchemaValue: sixty.IntValue(9), IsValid: true},
		sixty.ArrayValue{sixty.BitValue(true), sixty.BitValue(false)},
		sixty.ArrayValue{sixty.BitValue(true), sixty.BitValue(false)},
		sixty.ProductValue{sixty.IntValue(1), sixty.IntValue(-2)},
		sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1},
		sixty.ListValue{sixty.IntValue(1), sixty.IntValue(-2)},
		sixty.StringValue("hello"),
	})
	require.NoError(t, err)

	var dec allTypes
	err = stream.Decode(schema, enc, &dec)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

func TestDecodeIntAndEnumIntoAny(t *testing.T) {
	t.Parallel()

	type s struct {
		Int     any
		Uint    any
		Enum    any
		EnumInt any
		Str     any
		Bit     any
		Rune    any
	}

	schema := sixty.Product{
		{"f_int", sixty.Int{}},
		{"f_uint", sixty.Uint{}},
		{"f_enum", sixty.Enum{Values: []string{"a", "b"}}},
		{"f_enum_int", sixty.EnumInt{Min: 0, Max: 5}},
		{"f_str", sixty.String{}},
		{"f_bit", sixty.Bit{}},
		{"f_rune", sixty.Rune{}},
	}

	data, err := stream.Encode(schema, sixty.ProductValue{
		sixty.IntValue(-1),
		sixty.UintValue(1),
		sixty.EnumValue(1),
		sixty.EnumIntValue(3),
		sixty.StringValue("x"),
		sixty.BitValue(true),
		sixty.RuneValue('a'),
	})
	require.NoError(t, err)

	var got s
	err = stream.Decode(schema, data, &got)
	require.NoError(t, err)
	require.Equal(t, int64(-1), got.Int)
	require.Equal(t, uint64(1), got.Uint)
	require.Equal(t, uint(1), got.Enum)
	require.Equal(t, int64(3), got.EnumInt)
	require.Equal(t, "x", got.Str)
	require.Equal(t, true, got.Bit)
	require.Equal(t, 'a', got.Rune)
}

func TestDecodeWrongKindsPanic(t *testing.T) {
	t.Parallel()

	// string into non-string field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.String{}}},
			[]byte{0x01, 0x61}, &got)
	})
	// bit into non-bool field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.Bit{}}},
			[]byte{0x01}, &got)
	})
	// product into non-struct field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.Product{{"y", sixty.Int{}}}}},
			[]byte{0x02}, &got)
	})
	// coproduct into non-interface field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.Coproduct{{"y", sixty.Int{}}}}},
			[]byte{0x00, 0x02}, &got)
	})
	// list into non-slice field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.List{Schema: sixty.Int{}}}},
			[]byte{0x00}, &got)
	})
	// array into non-slice field
	require.Panics(t, func() {
		var got struct{ X int64 }
		_ = stream.Decode(sixty.Product{{"x", sixty.Array{Schema: sixty.Int{}, Count: 1}}},
			[]byte{0x02}, &got)
	})
}

func TestDecodeInvalidValuesPanic(t *testing.T) {
	t.Parallel()

	// bit marker 2
	require.Panics(t, func() {
		var got bool
		_ = stream.Decode(sixty.Bit{}, []byte{0x02}, &got)
	})
	// enum index out of range
	require.Panics(t, func() {
		var got uint8
		_ = stream.Decode(sixty.Enum{Values: []string{"a"}}, []byte{0x01}, &got)
	})
	// optional marker 2
	require.Panics(t, func() {
		var got *int64
		_ = stream.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x02}, &got)
	})
}
