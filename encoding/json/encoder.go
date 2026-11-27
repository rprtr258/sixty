package json

import (
	"bytes"
	"strconv"

	"github.com/rprtr258/sixty"
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
		if vv {
			e.Buffer.WriteString("true")
		} else {
			e.Buffer.WriteString("false")
		}
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
		for i, vvv := range vv {
			e.Buffer.WriteRune('"')
			e.Buffer.WriteString(tt[i].Name)
			e.Buffer.WriteString(`":`)
			e.Write(tt[i].Schema, vvv)
			if i < len(vv)-1 {
				e.Buffer.WriteRune(',')
			}
		}
		e.Buffer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.Buffer.WriteString(`{"`)
		e.Buffer.WriteString(strconv.FormatUint(uint64(vv.Variant), 10))
		e.Buffer.WriteString(`":`)
		e.Write(tt[vv.Variant].Schema, vv.SchemaValue)
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
