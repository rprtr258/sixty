package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
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
	require.NoError(t, err)
	require.Equal(t, int64(42), got)
}

// customSchema is not one of the sixty root schema types, driving the
// unsupported-schema error paths.
type customSchema struct{}

func (customSchema) Kind() sixty.Kind { return sixty.BitKind }

func TestDecodeIntoAny(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		value  sixty.SchemaValue
		want   any
	}{
		{sixty.String{}, sixty.StringValue("hi"), "hi"},
		{sixty.Bit{}, sixty.BitValue(true), true},
		{sixty.Int{}, sixty.IntValue(-3), int64(-3)},
		{sixty.Uint{}, sixty.UintValue(3), uint64(3)},
		{sixty.Rune{}, sixty.RuneValue('ж'), 'ж'},
		{sixty.Enum{Values: []string{"a", "b"}}, sixty.EnumValue(1), uint(1)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(-3), int64(-3)},
		// list/array/product require concrete slice/struct targets, only
		// optional and coproduct accept interface targets
		{sixty.Optional{Schema: sixty.Int{}}, sixty.OptionalValue{SchemaValue: sixty.IntValue(1), IsValid: true}, int64(1)},
	}
	for _, test := range testCases {
		var got any
		err := common.Decode(test.schema, test.value, &got)
		require.NoError(t, err, "schema %T", test.schema)
		require.Equal(t, test.want, got, "schema %T", test.schema)
	}
}

func TestEncodeUnsupportedSchema(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(customSchema{}, true)
	require.EqualError(t, err, "unsupported schema common_test.customSchema")
}

func TestDecodeUnsupportedSchema(t *testing.T) {
	t.Parallel()

	var got any
	err := common.Decode(customSchema{}, sixty.BitValue(true), &got)
	require.EqualError(t, err, "unsupported schema common_test.customSchema")
}

func TestDecodeNonPointerTarget(t *testing.T) {
	t.Parallel()

	var got int64
	err := common.Decode(sixty.Int{}, sixty.IntValue(1), got)
	require.EqualError(t, err, "decode target must be a pointer, got int64")
}

func TestDecodeValueMismatches(t *testing.T) {
	t.Parallel()

	// right schema kinds, wrong SchemaValue types
	testCases := []struct {
		schema sixty.Schema
		value  sixty.SchemaValue
		reason string
	}{
		{sixty.String{}, sixty.IntValue(1), "expected string value, got sixty.IntValue"},
		{sixty.Bit{}, sixty.IntValue(1), "expected bit value, got sixty.IntValue"},
		{sixty.Int{}, sixty.StringValue("x"), "expected int value, got sixty.StringValue"},
		{sixty.Uint{}, sixty.IntValue(1), "expected uint value, got sixty.IntValue"},
		{sixty.Rune{}, sixty.IntValue(1), "expected rune value, got sixty.IntValue"},
		{sixty.Enum{}, sixty.IntValue(1), "expected enum value, got sixty.IntValue"},
		{sixty.EnumInt{}, sixty.IntValue(1), "expected enum int value, got sixty.IntValue"},
		{sixty.Product{}, sixty.IntValue(1), "expected product value, got sixty.IntValue"},
		{sixty.Coproduct{}, sixty.IntValue(1), "expected coproduct value, got sixty.IntValue"},
		{sixty.List{}, sixty.IntValue(1), "expected list value, got sixty.IntValue"},
		{sixty.Array{}, sixty.IntValue(1), "expected array value, got sixty.IntValue"},
		{sixty.Optional{}, sixty.IntValue(1), "expected optional value, got sixty.IntValue"},
	}
	for _, test := range testCases {
		var got any
		err := common.Decode(test.schema, test.value, &got)
		require.EqualError(t, err, test.reason, "schema %T", test.schema)
	}
}

func TestDecodeKindMismatches(t *testing.T) {
	t.Parallel()

	// right SchemaValue types, wrong Go target kinds
	testCases := []struct {
		schema sixty.Schema
		value  sixty.SchemaValue
		target any
		reason string
	}{
		{sixty.String{}, sixty.StringValue("x"), new(int64), "expected string or any, got int64"},
		{sixty.Bit{}, sixty.BitValue(true), new(int64), "expected bool or any, got int64"},
		{sixty.Int{}, sixty.IntValue(1), new(uint64), "expected signed integer or any, got uint64"},
		{sixty.Uint{}, sixty.UintValue(1), new(int64), "expected unsigned integer or any, got int64"},
		{sixty.Rune{}, sixty.RuneValue('a'), new(int64), "expected rune or any, got int64"},
		{sixty.Enum{}, sixty.EnumValue(0), new(int64), "expected unsigned integer or any, got int64"},
		{sixty.EnumInt{}, sixty.EnumIntValue(0), new(uint64), "expected signed integer or any, got uint64"},
		{sixty.Product{}, sixty.ProductValue{}, new(int64), "expected struct, got int64"},
		{sixty.Coproduct{}, sixty.CoproductValue{}, new(int64), "expected interface, got int64"},
		{sixty.List{}, sixty.ListValue{}, new(int64), "expected slice, got int64"},
		{sixty.Array{Count: 1}, sixty.ArrayValue{sixty.IntValue(1)}, new(int64), "expected slice or array, got int64"},
	}
	for _, test := range testCases {
		err := common.Decode(test.schema, test.value, test.target)
		require.EqualError(t, err, test.reason, "schema %T", test.schema)
	}
}

func TestDecodeProductFieldCounts(t *testing.T) {
	t.Parallel()

	type pair struct{ X, Y int64 }

	var got pair
	err := common.Decode(
		sixty.Product{{"x", sixty.Int{}}, {"y", sixty.Int{}}},
		sixty.ProductValue{sixty.IntValue(1)},
		&got,
	)
	require.EqualError(t, err, "expected product value with 2 fields, got 1")

	single := struct{ X int64 }{}
	err = common.Decode(
		sixty.Product{{"x", sixty.Int{}}, {"y", sixty.Int{}}},
		sixty.ProductValue{sixty.IntValue(1), sixty.IntValue(2)},
		&single,
	)
	require.EqualError(t, err, "expected struct with 2 fields, got 1")
}

func TestDecodeProductNestedPathError(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{{"outer", sixty.Product{{"inner", sixty.Bit{}}}}}
	value := sixty.ProductValue{
		sixty.ProductValue{sixty.IntValue(1)},
	}
	var got struct{ Outer struct{ Inner bool } }
	err := common.Decode(schema, value, &got)
	require.EqualError(t, err, `field "outer": field "inner": expected bit value, got sixty.IntValue`)
}

func TestDecodeCoproductVariantOutOfRange(t *testing.T) {
	t.Parallel()

	var got any
	err := common.Decode(
		sixty.Coproduct{{"a", sixty.Int{}}},
		sixty.CoproductValue{Variant: 1, SchemaValue: sixty.IntValue(1)},
		&got,
	)
	require.EqualError(t, err, "coproduct variant 1 out of range 1")
}

func TestDecodeArrayLengthMismatch(t *testing.T) {
	t.Parallel()

	var got []int64
	err := common.Decode(
		sixty.Array{Schema: sixty.Int{}, Count: 2},
		sixty.ArrayValue{sixty.IntValue(1)},
		&got,
	)
	require.EqualError(t, err, "expected array value of 2 elements, got 1")
}

func TestDecodeArrayElementPathError(t *testing.T) {
	t.Parallel()

	var got []bool
	err := common.Decode(
		sixty.Array{Schema: sixty.Bit{}, Count: 2},
		sixty.ArrayValue{sixty.BitValue(true), sixty.IntValue(0)},
		&got,
	)
	require.EqualError(t, err, "element 1: expected bit value, got sixty.IntValue")
}

func TestDecodeListElementPathError(t *testing.T) {
	t.Parallel()

	var got []bool
	err := common.Decode(
		sixty.List{Schema: sixty.Bit{}},
		sixty.ListValue{sixty.IntValue(0)},
		&got,
	)
	require.EqualError(t, err, "element 0: expected bit value, got sixty.IntValue")
}

func TestDecodeOptional(t *testing.T) {
	t.Parallel()

	schema := sixty.Optional{Schema: sixty.Int{}}

	// absent zeroes the target
	var got int64
	err := common.Decode(schema, sixty.OptionalValue{IsValid: false}, &got)
	require.NoError(t, err)
	require.Equal(t, int64(0), got)

	// present into non-pointer reuses the target
	got = 5
	err = common.Decode(schema,
		sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true}, &got)
	require.NoError(t, err)
	require.Equal(t, int64(7), got)

	// present inner error propagates unwrapped
	optSchema := sixty.Optional{Schema: sixty.Bit{}}
	err = common.Decode(optSchema,
		sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true}, new(*bool))
	require.EqualError(t, err, "expected bit value, got sixty.IntValue")
}

func TestEncodeKindMismatches(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		value  any
		reason string
	}{
		{sixty.String{}, struct{ X int64 }{3}, `field "x": expected string, got int64`},
		{sixty.Bit{}, struct{ X string }{"yes"}, `field "x": expected bool, got string`},
		{sixty.Int{}, struct{ X uint64 }{3}, `field "x": expected signed integer, got uint64`},
		{sixty.Uint{}, struct{ X int64 }{3}, `field "x": expected unsigned integer, got int64`},
		{sixty.Rune{}, struct{ X int64 }{3}, `field "x": expected rune, got int64`},
		{sixty.Enum{Values: []string{"a"}}, struct{ X int64 }{0}, `field "x": expected unsigned integer, got int64`},
		{sixty.Enum{Values: []string{"a"}}, struct{ X uint64 }{1}, `field "x": enum index 1 out of range 1`},
		{sixty.EnumInt{Min: 0, Max: 5}, struct{ X uint64 }{3}, `field "x": expected signed integer, got uint64`},
		{sixty.EnumInt{Min: 0, Max: 5}, struct{ X int64 }{6}, `field "x": value 6 out of enum int range [0..5]`},
		// containers encode the bare value, so mismatch messages have no field prefix
		{sixty.Product{{"x", sixty.Int{}}}, 42, `expected struct, got int`},
		{sixty.List{Schema: sixty.Int{}}, 42, `expected slice, got int`},
		{sixty.Array{Schema: sixty.Int{}, Count: 2}, 42, `expected array, got int`},
	}
	for _, test := range testCases {
		schema := test.schema
		if _, isContainer := test.value.(int); isContainer {
			_, err := common.Encode(schema, test.value)
			require.EqualError(t, err, test.reason, "schema %T", test.schema)
		} else {
			_, err := common.Encode(sixty.Product{{"x", schema}}, test.value)
			require.EqualError(t, err, test.reason, "schema %T", test.schema)
		}
	}
}

func TestEncodeProductFieldCountMismatch(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(sixty.Product{{"x", sixty.Int{}}}, struct{ X, Y int64 }{1, 2})
	require.EqualError(t, err, "expected struct with 1 fields, got 2")
}

func TestEncodeProductNestedPathError(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{{"outer", sixty.Product{{"inner", sixty.Bit{}}}}}
	_, err := common.Encode(schema, struct{ Outer struct{ Inner string } }{})
	require.EqualError(t, err, `field "outer": field "inner": expected bool, got string`)

	var pathErr common.PathError
	require.ErrorAs(t, err, &pathErr)
	require.Equal(t, `field "outer"`, pathErr.Path)
	var mismatch common.MismatchError
	require.ErrorAs(t, pathErr.Unwrap(), &mismatch)
}

func TestEncodeCoproductFirstMatchWins(t *testing.T) {
	t.Parallel()

	schema := sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}

	// bare top-level values are matched against the variants directly
	got, err := common.Encode(schema, int64(42))
	require.NoError(t, err)
	require.Equal(t, sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0}, got)

	got, err = common.Encode(schema, "hi")
	require.NoError(t, err)
	require.Equal(t, sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1}, got)

	got, err = common.Encode(sixty.Product{{"v", schema}}, struct{ V any }{int64(42)})
	require.NoError(t, err)
	require.Equal(t,
		sixty.ProductValue{sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0}},
		got,
	)

	got, err = common.Encode(sixty.Product{{"v", schema}}, struct{ V any }{"hi"})
	require.NoError(t, err)
	require.Equal(t,
		sixty.ProductValue{sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1}},
		got,
	)
}

func TestEncodeCoproductNoMatch(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(
		sixty.Product{{"v", sixty.Coproduct{{"x", sixty.String{}}}}},
		struct{ V any }{int64(42)},
	)
	require.EqualError(t, err, `field "v": no coproduct variant matches value of kind int64`)
}

func TestEncodeArrayLengthMismatch(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(sixty.Array{Schema: sixty.Int{}, Count: 2}, []int64{1})
	require.EqualError(t, err, "expected array of 2 elements, got 1")
}

func TestEncodeArrayElementPathError(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(sixty.Array{Schema: sixty.Bit{}, Count: 2}, []string{"a", "b"})
	require.EqualError(t, err, "element 0: expected bool, got string")
}

func TestEncodeListElementPathError(t *testing.T) {
	t.Parallel()

	_, err := common.Encode(sixty.List{Schema: sixty.Bit{}}, []string{"a"})
	require.EqualError(t, err, "element 0: expected bool, got string")
}

func TestDecodeArrayIntoFixedArray(t *testing.T) {
	t.Parallel()

	var got [2]bool
	err := common.Decode(
		sixty.Array{Schema: sixty.Bit{}, Count: 2},
		sixty.ArrayValue{sixty.BitValue(true), sixty.BitValue(false)},
		&got,
	)
	require.NoError(t, err)
	require.Equal(t, [2]bool{true, false}, got)

	var gotSlice []bool
	err = common.Decode(
		sixty.Array{Schema: sixty.Bit{}, Count: 2},
		sixty.ArrayValue{sixty.BitValue(true), sixty.BitValue(false)},
		&gotSlice,
	)
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, gotSlice)
}

func TestEncodeOptional(t *testing.T) {
	t.Parallel()

	schema := sixty.Optional{Schema: sixty.Int{}}

	opt := int64(7)
	got, err := common.Encode(schema, &opt)
	require.NoError(t, err)
	require.Equal(t, sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true}, got)

	got, err = common.Encode(schema, nil)
	require.NoError(t, err)
	require.Equal(t, sixty.OptionalValue{IsValid: false}, got)

	_, err = common.Encode(schema, &struct{}{})
	require.ErrorContains(t, err, "expected signed integer, got struct")
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
		{"i", sixty.Int{}},
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
	for _, data := range []sample{{
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
		enc, err := common.Encode(schema, data)
		require.NoError(t, err)
		encBytes, err := bytes.Encode(schema, enc)
		require.NoError(t, err)
		decValue, err := bytes.Decode(schema, encBytes)
		require.NoError(t, err)
		var dec sample
		err = common.Decode(schema, decValue, &dec)
		require.NoError(t, err)
		require.Equal(t, dec, data)
	}
}
