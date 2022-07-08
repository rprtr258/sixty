package bytesenc

import (
	sixty "github.com/rprtr258/sixty/pkg"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
// TODO: compress list/array of bits
// TODO: streaming encoder/decoder, with size of element for lists (if it is not known)
type ByteEncoder struct {
	// TODO: make unexported
	Buffer interface {
		WriteByte(byte) error
		WriteString(string) (int, error)
		WriteRune(rune) (int, error)
	}
}

func (v *ByteEncoder) writeUint(x uint) {
	// TODO: implement
	switch {
	case x < 128:
		v.Buffer.WriteByte(byte(x))
	}
}
func (v *ByteEncoder) writeInt(x int) {
	if x > 0 {
		v.writeUint(uint(x) * 2)
	} else {
		// TODO: fix and test overflow/steal from protobuf
		v.writeUint(uint(-2*x) - 1)
	}
}

func assertSchemas[S sixty.Schema, V sixty.SchemaValue](t sixty.Schema, v sixty.SchemaValue) (tt S, vv V, ok bool) {
	tt, okSchema := t.(S)
	vv, okValue := v.(V)
	ok = okSchema && okValue
	return
}

func (e *ByteEncoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := assertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writeUint(uint(len(vv)))
		e.Buffer.WriteString(string(vv))
	} else if _, vv, ok := assertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		if vv {
			e.Buffer.WriteByte(1)
		} else {
			e.Buffer.WriteByte(0)
		}
	} else if _, vv, ok := assertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writeInt(int(vv))
	} else if _, vv, ok := assertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writeUint(uint(vv))
	} else if _, vv, ok := assertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.Buffer.WriteRune(rune(vv))
	} else if _, vv, ok := assertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.writeUint(uint(vv))
	} else if tt, vv, ok := assertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.writeUint(uint(int(vv) - tt.Min))
	} else if tt, vv, ok := assertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i], x)
		}
	} else if tt, vv, ok := assertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.writeUint(vv.Variant)
		e.Write(tt[vv.Variant], vv.SchemaValue)
	} else if tt, vv, ok := assertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		e.writeUint(uint(len(vv)))
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := assertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := assertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.writeUint(0)
		} else {
			e.writeUint(1)
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else if tt, vv, ok := assertSchemas[sixty.Alias, sixty.AliasValue](t, v); ok {
		e.Write(tt.Schema, vv.Value)
	} else {
		panic("something unimplemented")
	}
}
