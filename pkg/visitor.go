package sixty

// // TODO: change to kinds
// func checkTypeValue(t Schema, val SchemaValue, kind Kind) bool {
// 	return t.Kind() == kind && val.Kind() == kind
// }

// func VisitValue(t Schema, val SchemaValue, visitor interface {
// 	VisitBit(Bit, BitValue)
// 	VisitInt(Int, IntValue)
// 	VisitUint(Uint, UintValue)
// 	VisitRune(Rune, RuneValue)
// 	VisitEnum(Enum, EnumValue)
// 	VisitEnumInt(EnumInt, EnumIntValue)
// 	VisitProduct(Product, ProductValue)
// 	VisitCoproduct(Coproduct, CoproductValue)
// 	VisitList(List, ListValue)
// 	VisitArray(Array, ArrayValue)
// 	VisitOptional(Optional, OptionalValue)
// 	VisitAlias(Alias, AliasValue)
// 	VisitString(string)
// }) {
// 	if t == (String{}) {
// 		// TODO: check val is string
// 		s := val.(AliasValue)
// 		var b strings.Builder
// 		for _, c := range s.Value.(ListValue) {
// 			cc := c.(RuneValue)
// 			b.WriteRune(rune(cc))
// 		}
// 		visitor.VisitString(b.String())
// 		return
// 	}
// 	if t.Kind() != val.Kind() {
// 		// TODO: error: wrong schema
// 		panic(fmt.Sprintf("%T(%[1]v) is not %T(%[2]v) value", val, t))
// 	}
// 	switch t.Kind() {
// 	case BitKind:
// 		visitor.VisitBit(t.(Bit), val.(BitValue))
// 	case IntKind:
// 		visitor.VisitInt(t.(Int), val.(IntValue))
// 	case UintKind:
// 		visitor.VisitUint(t.(Uint), val.(UintValue))
// 	case RuneKind:
// 		visitor.VisitRune(t.(Rune), val.(RuneValue))
// 	case EnumKind:
// 		visitor.VisitEnum(t.(Enum), val.(EnumValue))
// 	case EnumIntKind:
// 		visitor.VisitEnumInt(t.(EnumInt), val.(EnumIntValue))
// 	case ProductKind:
// 		visitor.VisitProduct(t.(Product), val.(ProductValue))
// 	case CoproductKind:
// 		visitor.VisitCoproduct(t.(Coproduct), val.(CoproductValue))
// 	case ListKind:
// 		visitor.VisitList(t.(List), val.(ListValue))
// 	case ArrayKind:
// 		visitor.VisitArray(t.(Array), val.(ArrayValue))
// 	case OptionalKind:
// 		visitor.VisitOptional(t.(Optional), val.(OptionalValue))
// 	case AliasKind:
// 		visitor.VisitAlias(t.(Alias), val.(AliasValue))
// 	default:
// 		panic(fmt.Sprintf("unimplemented kind %T", t))
// 	}
// }

// func VisitValue2(t Schema, visitor interface {
// 	VisitBit(Bit) BitValue
// 	VisitInt(Int) IntValue
// 	VisitUint(Uint) UintValue
// 	VisitRune(Rune) RuneValue
// 	VisitEnum(Enum) EnumValue
// 	VisitEnumInt(EnumInt) EnumIntValue
// 	VisitProduct(Product) ProductValue
// 	VisitCoproduct(Coproduct) CoproductValue
// 	VisitList(List) ListValue
// 	VisitArray(Array) ArrayValue
// 	VisitOptional(Optional) OptionalValue
// 	VisitAlias(Alias) AliasValue
// 	VisitString() string
// }) SchemaValue {
// 	if t == (String{}) {
// 		s := visitor.VisitString()
// 		runes := make([]SchemaValue, 0, len(s))
// 		for _, c := range s {
// 			runes = append(runes, RuneValue(c))
// 		}
// 		return AliasValue{
// 			Name:  string(t.(String).PredefKind()),
// 			Value: ListValue(runes),
// 		}
// 	}

// 	switch tt := t.(type) {
// 	case Bit:
// 		return visitor.VisitBit(tt)
// 	case Int:
// 		return visitor.VisitInt(tt)
// 	case Uint:
// 		return visitor.VisitUint(tt)
// 	case Rune:
// 		return visitor.VisitRune(tt)
// 	case Enum:
// 		return visitor.VisitEnum(tt)
// 	case EnumInt:
// 		return visitor.VisitEnumInt(tt)
// 	case Product:
// 		return visitor.VisitProduct(tt)
// 	case Coproduct:
// 		return visitor.VisitCoproduct(tt)
// 	case List:
// 		return visitor.VisitList(tt)
// 	case Array:
// 		return visitor.VisitArray(tt)
// 	case Optional:
// 		return visitor.VisitOptional(tt)
// 	case Alias:
// 		return visitor.VisitAlias(tt)
// 	default:
// 		// TODO: wtf
// 		panic("wtf")
// 	}
// }
