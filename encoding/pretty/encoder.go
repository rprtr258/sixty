package pretty

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal/stringify"
)

type Encoder struct {
	writer sixty.Writer
	indent int
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	indent := strings.Repeat("  ", e.indent)
	e.writer.WriteString(indent)
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writer.WriteRune('"')
		e.writer.WriteString(strings.ReplaceAll(string(vv), "\"", "\"\""))
		e.writer.WriteRune('"')
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		e.writer.WriteString(stringify.BoolToText(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writer.WriteString(strconv.FormatInt(int64(vv), 10))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writer.WriteString(strconv.FormatUint(uint64(vv), 10))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.writer.WriteRune('\'')
		if vv == '\'' {
			e.writer.WriteRune('\\')
		}
		e.writer.WriteRune(rune(vv))
		e.writer.WriteRune('\'')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.writer.WriteString(tt.Values[vv])
	} else if _, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.writer.WriteString(strconv.FormatInt(int64(vv), 10))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		e.writer.WriteString("{\n")
		e.indent++
		for i, x := range vv {
			e.Write(tt[i].Schema, x)
			e.writer.WriteString(",\n")
		}
		e.indent--
		e.writer.WriteString(indent)
		e.writer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.writer.WriteString("{\n")
		e.indent++
		e.writer.WriteString(strconv.FormatUint(uint64(vv.Variant), 10))
		e.writer.WriteString(`":`)
		e.Write(tt[vv.Variant].Schema, vv.SchemaValue)
		e.indent--
		e.writer.WriteString(indent)
		e.writer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		e.writer.WriteString("[\n")
		e.indent++
		for i := range vv {
			e.Write(tt.Schema, vv[i])
			e.writer.WriteString(",\n")
		}
		e.indent--
		e.writer.WriteString(indent)
		e.writer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		e.writer.WriteString("[\n")
		for i := uint(0); i < tt.Count; i++ {
			e.Write(tt.Schema, vv[i])
			e.writer.WriteString(",\n")
		}
		e.writer.WriteString(indent)
		e.writer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.writer.WriteString("null")
		} else {
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else {
		panic("something unimplemented")
	}
}

func Encode(t sixty.Schema, v sixty.SchemaValue) (string, error) {
	buffer := bytes.Buffer{}
	encoder := Encoder{writer: sixty.NewWriter(&buffer)}
	encoder.Write(t, v)
	return buffer.String(), encoder.writer.Error
}
