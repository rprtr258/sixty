package json

import (
	"bytes"
	"strconv"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal/stringify"
)

type Encoder struct {
	Buffer sixty.Writer
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.Buffer.WriteRune('"')
		e.Buffer.WriteString(string(vv))
		e.Buffer.WriteRune('"')
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		e.Buffer.WriteString(stringify.BoolToText(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.Buffer.WriteString(strconv.FormatInt(int64(vv), 10))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.Buffer.WriteString(strconv.FormatUint(uint64(vv), 10))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.Buffer.WriteRune('"')
		e.Buffer.WriteRune(rune(vv))
		e.Buffer.WriteRune('"')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.Buffer.WriteString(tt.Values[vv])
	} else if _, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.Buffer.WriteString(strconv.FormatInt(int64(vv), 10))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		e.Buffer.WriteRune('{')
		for i, x := range vv {
			// TODO: remove
			// if ttt, vvv, ok := sixty.AssertSchemas[sixty.Alias, sixty.AliasValue](tt[i], vv); ok {
			// 	e.Buffer.WriteRune('"')
			// 	e.Buffer.WriteString(ttt.Name)
			// 	e.Buffer.WriteString(`":`)
			// 	e.Write(ttt.Schema, vvv.Value)
			// } else {
			e.Buffer.WriteRune('"')
			e.Buffer.WriteString(strconv.FormatInt(int64(i), 10))
			e.Buffer.WriteString(`":`)
			e.Write(tt[i], x)
			// }
			if i < len(vv)-1 {
				e.Buffer.WriteRune(',')
			}
		}
		e.Buffer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.Buffer.WriteString(`{"`)
		e.Buffer.WriteString(strconv.FormatUint(uint64(vv.Variant), 10))
		e.Buffer.WriteString(`":`)
		e.Write(tt[vv.Variant], vv.SchemaValue)
		e.Buffer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		e.Buffer.WriteRune('[')
		i := 0
	iter_list:
		e.Write(tt.Schema, vv[i])
		if i < len(vv)-1 {
			i++
			e.Buffer.WriteRune(',')
			goto iter_list
		}
		e.Buffer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		e.Buffer.WriteRune('[')
		i := uint(0)
	iter_array:
		e.Write(tt.Schema, vv[i])
		if i < tt.Count-1 {
			i++
			e.Buffer.WriteRune(',')
			goto iter_array
		}
		e.Buffer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.Buffer.WriteString("null")
		} else {
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else {
		panic("something unimplemented")
	}
}

func Encode(t sixty.Schema, v sixty.SchemaValue) (string, error) {
	buffer := bytes.Buffer{}
	encoder := Encoder{Buffer: sixty.NewWriter(&buffer)}
	encoder.Write(t, v)
	return buffer.String(), encoder.Buffer.Error
}
