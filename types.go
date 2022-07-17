package sixty

type Kind uint

const (
	BitKind Kind = iota
	IntKind
	UintKind
	RuneKind
	EnumKind
	EnumIntKind
	ProductKind
	CoproductKind
	ListKind
	ArrayKind
	OptionalKind
	AliasKind
)

type PredefKind string

const (
	StringKind PredefKind = ":string"
	// Date
	// Datetime
	//    millis,micros,nanos
	// []bit
	// [n]bit
)

type Schema interface {
	Kind() Kind
	// TODO: get value?
}

type PredefSchema interface {
	PredefKind() PredefKind
	Schema() Schema
	// TODO: get value?
}

type SchemaValue interface {
	isSchemaValue()
	Kind() Kind
}

type Bit struct{}

func (Bit) Kind() Kind { return BitKind }

type BitValue bool

func (BitValue) isSchemaValue() {}
func (BitValue) Kind() Kind     { return BitKind }

type Int struct{}

func (Int) Kind() Kind { return IntKind }

type IntValue int64

func (IntValue) isSchemaValue() {}
func (IntValue) Kind() Kind     { return IntKind }

type Uint struct{}

func (Uint) Kind() Kind { return UintKind }

type UintValue uint64

func (UintValue) isSchemaValue() {}
func (UintValue) Kind() Kind     { return UintKind }

type Rune struct{}

func (Rune) Kind() Kind { return RuneKind }

type RuneValue rune

func (RuneValue) isSchemaValue() {}
func (RuneValue) Kind() Kind     { return RuneKind }

type Enum struct {
	Values []string
}

func (Enum) Kind() Kind { return EnumKind }

type EnumValue uint

func (EnumValue) isSchemaValue() {}
func (EnumValue) Kind() Kind     { return EnumKind }

type EnumInt struct {
	Min int64
	Max int64
}

func (EnumInt) Kind() Kind { return EnumIntKind }

type EnumIntValue int64

func (EnumIntValue) isSchemaValue() {}
func (EnumIntValue) Kind() Kind     { return EnumIntKind }

type Product []Schema

func (Product) Kind() Kind { return ProductKind }

type ProductValue []SchemaValue

func (ProductValue) isSchemaValue() {}
func (ProductValue) Kind() Kind     { return ProductKind }

type Coproduct []Schema

func (Coproduct) Kind() Kind { return CoproductKind }

type CoproductValue struct {
	SchemaValue
	Variant uint
}

func (CoproductValue) isSchemaValue() {}
func (CoproductValue) Kind() Kind     { return CoproductKind }

type List struct {
	Schema
}

func (List) Kind() Kind { return ListKind }

type ListValue []SchemaValue

func (ListValue) isSchemaValue() {}
func (ListValue) Kind() Kind     { return ListKind }

type Array struct {
	Schema
	Count uint
}

func (Array) Kind() Kind { return ArrayKind }

type ArrayValue []SchemaValue

func (ArrayValue) isSchemaValue() {}
func (ArrayValue) Kind() Kind     { return ArrayKind }

type Optional struct {
	Schema
}

func (Optional) Kind() Kind { return OptionalKind }

type OptionalValue struct {
	SchemaValue
	IsValid bool
}

func (OptionalValue) isSchemaValue() {}
func (OptionalValue) Kind() Kind     { return OptionalKind }

// TODO: alias in alias???
type Alias struct {
	Schema
	Name string
}

func (Alias) Kind() Kind { return AliasKind }

type AliasValue struct {
	Name  string
	Value SchemaValue
}

func (AliasValue) isSchemaValue() {}
func (AliasValue) Kind() Kind     { return AliasKind }

type String struct{}

func (s String) Kind() Kind           { return s.Schema().Kind() }
func (String) PredefKind() PredefKind { return StringKind }
func (String) Schema() Schema         { return List{Schema: Rune{}} }

type StringValue string

func (StringValue) isSchemaValue() {}
func (StringValue) Kind() Kind     { return ListKind }

// type TypeWithMetadata struct {
// 	Type
// 	KVMeta    map[string]string
// 	FlagsMeta map[string]struct{}
// }

// TODO: Metadata String
