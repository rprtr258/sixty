package stream

import (
	"bytes"

	"github.com/rprtr258/sixty"
)

type Encoder struct {
	writer sixty.Writer
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writer.WriteUint64(uint64(len(vv)))
		e.writer.WriteString(string(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		if vv {
			e.writer.WriteByte(1)
		} else {
			e.writer.WriteByte(0)
		}
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writer.WriteInt64(int64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.writer.WriteRune(rune(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.writer.WriteUint64(uint64(int64(vv) - tt.Min))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i], x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv.Variant))
		e.Write(tt[vv.Variant], vv.SchemaValue)
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		for _, x := range vv {
			// TODO: maybe encode as how many elements goes next for sure
			e.writer.WriteUint64(1)
			e.Write(tt.Schema, x)
		}
		e.writer.WriteUint64(0)
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.writer.WriteUint64(0)
		} else {
			e.writer.WriteUint64(1)
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Alias, sixty.AliasValue](t, v); ok {
		e.Write(tt.Schema, vv.Value)
	} else {
		panic("something unimplemented")
	}
}

func Encode(t sixty.Schema, v sixty.SchemaValue) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := Encoder{
		writer: sixty.NewWriter(&buffer),
	}
	encoder.Write(t, v)
	return buffer.Bytes(), encoder.writer.Error
}
