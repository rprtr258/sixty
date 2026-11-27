package stream_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/stream"
)

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

func TestDecodeWrongKinds(t *testing.T) {
	t.Parallel()

	// string into non-string field
	var strTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.String{}}},
			[]byte{0x01, 0x61}, &strTarget),
		"expected string, got int64")
	// bit into non-bool field
	var bitTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.Bit{}}},
			[]byte{0x01}, &bitTarget),
		"expected bool, got int64")
	// product into non-struct field
	var prodTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.Product{{"y", sixty.Int{}}}}},
			[]byte{0x02}, &prodTarget),
		"expected struct, got int64")
	// coproduct into non-interface field
	var copTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.Coproduct{{"y", sixty.Int{}}}}},
			[]byte{0x00, 0x02}, &copTarget),
		"expected interface, got int64")
	// list into non-slice field
	var listTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.List{Schema: sixty.Int{}}}},
			[]byte{0x00}, &listTarget),
		"expected slice, got int64")
	// array into non-slice field
	var arrTarget struct{ X int64 }
	require.ErrorContains(t,
		stream.Decode(sixty.Product{{"x", sixty.Array{Schema: sixty.Int{}, Count: 1}}},
			[]byte{0x02}, &arrTarget),
		"expected slice or array, got int64")
}

func TestDecodeInvalidValues(t *testing.T) {
	t.Parallel()

	// bit marker 2
	var bitGot bool
	require.ErrorContains(t,
		stream.Decode(sixty.Bit{}, []byte{0x02}, &bitGot),
		"expected bit which is either 0 or 1, found: 2")
	// enum index out of range
	var enumGot uint8
	require.ErrorContains(t,
		stream.Decode(sixty.Enum{Values: []string{"a"}}, []byte{0x01}, &enumGot),
		"found 1'th enum value while there are only 1 of them")
	// optional marker 2
	var optGot *int64
	require.ErrorContains(t,
		stream.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x02}, &optGot),
		"option variant must be either 0 or 1, but it is 2")
	// enum int out of range
	var enumIntGot int64
	require.ErrorContains(t,
		stream.Decode(sixty.EnumInt{Min: 0, Max: 1}, []byte{0x05}, &enumIntGot),
		"out of range")
	// coproduct variant out of range
	var copGot any
	require.ErrorContains(t,
		stream.Decode(sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}, []byte{0x03, 0x02}, &copGot),
		"found variant 3 while coproduct has only 2 of them")
	// invalid list marker
	var listGot []int64
	require.ErrorContains(t,
		stream.Decode(sixty.List{Schema: sixty.Int{}}, []byte{0x05}, &listGot),
		"expected list marker to be either 0 or 1, but it is 5")
}
