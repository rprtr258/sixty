package json

import (
	"bytes"
	"fmt"

	sixty "github.com/rprtr258/sixty/pkg"
)

type Encoder struct {
	Buffer *bytes.Buffer
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
		e.Buffer.WriteString(fmt.Sprint(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.Buffer.WriteString(fmt.Sprint(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.Buffer.WriteRune('"')
		e.Buffer.WriteRune(rune(vv))
		e.Buffer.WriteRune('"')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.Buffer.WriteString(tt.Values[vv])
	} else if _, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.Buffer.WriteString(fmt.Sprint(vv))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		e.Buffer.WriteRune('{')
		for i, x := range vv {
			// TODO: move out std:types as a separate case
			if ttt, vvv, ok := sixty.AssertSchemas[sixty.Alias, sixty.AliasValue](tt[i], vv); ok {
				e.Buffer.WriteRune('"')
				e.Buffer.WriteString(ttt.Name)
				e.Buffer.WriteString(`":`)
				e.Write(ttt.Schema, vvv.Value)
			} else {
				e.Buffer.WriteRune('"')
				e.Buffer.WriteString(fmt.Sprint(i))
				e.Buffer.WriteString(`":`)
				e.Write(tt[i], x)
			}
			if i < len(vv)-1 {
				e.Buffer.WriteRune(',')
			}
		}
		e.Buffer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.Buffer.WriteString(`{"`)
		e.Buffer.WriteString(fmt.Sprint(vv.Variant))
		e.Buffer.WriteString(`":`)
		e.Write(tt[vv.Variant], vv.SchemaValue)
		e.Buffer.WriteRune('}')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		e.Buffer.WriteRune('[')
		for i, x := range vv {
			e.Write(tt.Schema, x)
			if i < len(vv)-1 {
				e.Buffer.WriteRune(',')
			}
		}
		e.Buffer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		e.Buffer.WriteRune('[')
		for i, x := range vv {
			e.Write(tt.Schema, x)
			// TODO: iter while
			if uint(i) < tt.Count-1 {
				e.Buffer.WriteRune(',')
			}
		}
		e.Buffer.WriteRune(']')
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.Buffer.WriteString("null")
		} else {
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Alias, sixty.AliasValue](t, v); ok {
		e.Write(tt.Schema, vv.Value)
	} else {
		panic("something unimplemented")
	}
}
