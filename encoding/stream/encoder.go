package stream

import (
	"bytes"
	"unicode/utf8"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal/check"
)

type Encoder struct {
	writer sixty.Writer
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := check.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		// length prefix counts runes: ReadLengthString reads that many runes
		e.writer.WriteUint64(uint64(utf8.RuneCountInString(string(vv))))
		e.writer.WriteString(string(vv))
	} else if _, vv, ok := check.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		if vv {
			e.writer.WriteByte(1)
		} else {
			e.writer.WriteByte(0)
		}
	} else if _, vv, ok := check.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writer.WriteInt64(int64(vv))
	} else if _, vv, ok := check.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv))
	} else if _, vv, ok := check.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.writer.WriteRune(rune(vv))
	} else if _, vv, ok := check.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv))
	} else if tt, vv, ok := check.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.writer.WriteUint64(uint64(int64(vv) - tt.Min))
	} else if tt, vv, ok := check.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i].Schema, x)
		}
	} else if tt, vv, ok := check.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.writer.WriteUint64(uint64(vv.Variant))
		e.Write(tt[vv.Variant].Schema, vv.SchemaValue)
	} else if tt, vv, ok := check.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		for _, x := range vv {
			// 1/0 markers keep the format decodable element-by-element while
			// streaming, per the format design goals
			e.writer.WriteUint64(1)
			e.Write(tt.Schema, x)
		}
		e.writer.WriteUint64(0)
	} else if tt, vv, ok := check.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := check.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.writer.WriteUint64(0)
		} else {
			e.writer.WriteUint64(1)
			e.Write(tt.Schema, vv.SchemaValue)
		}
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
	encoder.writer.Flush()
	return buffer.Bytes(), encoder.writer.Error
}
