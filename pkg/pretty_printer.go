package sixty

import (
	"fmt"
	"strings"
)

type prettyPrinter struct {
	builder strings.Builder
	ident   uint
	aliases map[string]string
}

func NewPrettyPrinter(builder strings.Builder) prettyPrinter {
	return prettyPrinter{
		ident:   0,
		aliases: make(map[string]string),
		builder: builder,
	}
}

func (v *prettyPrinter) Write(t Schema) {
	// TODO: cache idents
	v.builder.WriteString(strings.Repeat(" ", int(v.ident*2)))
	switch tt := t.(type) {
	case String:
		v.builder.WriteString(string(StringKind))
	case Bit:
		v.builder.WriteString("bit")
	case Int:
		v.builder.WriteString("int")
	case Uint:
		v.builder.WriteString("uint")
	case Rune:
		v.builder.WriteString("rune")
	case Enum:
		v.builder.WriteRune('{')
		// TODO: iter while
		v.builder.WriteString(strings.Join(tt.Values, ", "))
		v.builder.WriteRune('}')
	case EnumInt:
		v.builder.WriteRune('(')
		v.builder.WriteString(fmt.Sprint(tt.Min))
		v.builder.WriteString("..")
		v.builder.WriteString(fmt.Sprint(tt.Max))
		v.builder.WriteRune(')')
	case Product:
		v.builder.WriteString("(\n")
		v.ident++
		for _, field := range tt {
			v.Write(field)
			v.builder.WriteString(",\n")
		}
		v.ident--
		v.builder.WriteString(")")
	case Coproduct:
		v.builder.WriteString("{\n")
		v.ident++
		for _, variant := range tt {
			v.Write(variant)
			v.builder.WriteString(",\n")
		}
		v.ident--
		v.builder.WriteString("}")
	case List:
		v.builder.WriteString("[]")
		tmp := v.ident
		v.ident = 0
		v.Write(tt.Schema)
		v.ident = tmp
	case Array:
		v.builder.WriteString(fmt.Sprintf("[%d]", tt.Count))
		tmp := v.ident
		v.ident = 0
		v.Write(tt.Schema)
		v.ident = tmp
	case Optional:
		v.builder.WriteString("?")
		tmp := v.ident
		v.ident = 0
		v.Write(tt.Schema)
		v.ident = tmp
	case Alias:
		v.builder.WriteString(tt.Name)
		if _, ok := v.aliases[tt.Name]; ok {
			return
		}
		pp := prettyPrinter{aliases: v.aliases}
		pp.Write(tt.Schema)
		v.aliases[tt.Name] = pp.String()
	default:
		panic("something unimplemented")
	}
}

func (v *prettyPrinter) String() string {
	var res strings.Builder
	for k, v := range v.aliases {
		res.WriteString(fmt.Sprintf("%s = %s\n\n", k, v))
	}
	res.WriteString(v.builder.String())
	return res.String()
}

func PrettyPrintSchema(t Schema) string {
	tsv := NewPrettyPrinter(strings.Builder{})
	tsv.Write(t)
	return tsv.String()
}
