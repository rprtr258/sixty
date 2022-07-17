# Sixty data format

## Design decisions
- Arbitary schema (unlike `CSV`). More precisely it would be ideal to have these types (`A`, `B`, `C`... are some arbitary types):
    - primitive types (at least): `int`, `float`, `string`, `bit`, `byte`
    - array of `T` of length `N`: `N` elements of type `T`
    - list of `T`: collection of elements of type `T` of indefinite length
    - optional `T`: either element of type `T` or `null` (nothing)
    - product of `A`, `B`, `C`...: like C struct with element of `A`, element of `B`, element of `C`, etc
    - sum of `A`, `B`, `C`...: like tagged union: element of this type is either `A`, or `B`, or `C`, etc.
- No or minimal excess data like dict keys in array elements in `JSON` or `<aboba>`, `</aboba>` stuff in `XML`.
- Schema definition is separated from specific encoding. So, there are two parts: schema definition rules and syntax, and encoding rules. Some usages might use binary encoding, some may use text, e.g. for presenting to humans.
- Support streaming decoding for list types. Really this depends on encoding, not schema definition.
- Schema might be provided either with data encoded or be predefined. If schema is coming with data (like headers row in `CSV`) it is used to parse data. If not.
- Back/forward compatibilities are not design goals. Speed encoding/decoding is not design goal. Encoding size is nearly design goal, not requiring encoding to be super efficient and compressed.
- Open and free for everyone.
- Simple to use. That means simple to understand, use and that I would write libs for languages I'm going to use the format in and cli tool to convert to/from that format. No dependencies on any language/system/ecosystem is required.

## Specification
Data format is divided into two parts: schema definition and data. Schema is about which types in which order are stored in data. Data is the data itself written using some encoding.

### Basic types
|name|description|
|--|--|
|`bit`|one bit of data|
|`int`, `uint`|arbitarily big, but usually bounded to 64 bits|
|`rune`|unicode char|

### Composite types
|name|syntax|description|
|--|--|--|
|product/struct/record/tuple|`(a=A, b=B, c=C)`||
|coproduct/choice/union|`[a=A, b=B, c=C]`|value of one of `A`,`B`,`C` types|
|enum|`{a, b, c}`|labels list of enum values|
|enum int|`[min..max]`|any value from `min` to `max` including|
|optional|`?A`|value of type `A` or nothing|
|list of|`[]A`|zero or more values of type `A`|
|array of|`[N]A`|`N` values of type `A`|

E.g. bitfield can be simulated in following way:
```
permission = (
  read = bit,
  write = bit,
  execute = bit,
)
full_permission = (
  owner = permission,
  group = permission,
  other = permission,
)
```

### Logical types
Logical types are types that can be represented using basic and composite types, but may have special datatype in language to be decoded to. So they are separated from their inner representation, which you can still use if you want.
|name|type|description|
|--|--|--|
|`:i8`, `:u8`|`[8]bit`||
|`:i16`, `:u16`|`[16]bit`||
|`:i32`, `:u32`|`[32]bit`||
|`:i64`, `:u64`|`[64]bit`||
|`:f16`|`(sign = bit, exponent = [5]bit, mantissa = [10]bit)`||
|`:f32`|`(sign = bit, exponent = [8]bit, mantissa = [23]bit)`||
|`:f64`|`(sign = bit, exponent = [11]bit, mantissa = [52]bit)`||
|`:string`|`[]rune`||
|`:datetime` ???|`uint`|unix timestamp with nanoseconds|
|`:bool` ???|`bit` or `{TRUE, FALSE}`||

### Alias
Types can be aliased to use later.
```
identifier ::= [A-Za-z_][A-Za-z0-9_-]*
```
Some examples:
```
month = [0..11]
week_day = {Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}
// week_day might also be described as [0..6] but then it is up to application to deduce what is first day of the week: monday or sunday
date = (
  day = u8,
  month = month,
  year = year,
)
second = [0..59]
minute = [0..59]
hour = [0..23]
time = (second, minute, hour)
// really it is just int - number of seconds since Unix epoch
datetime = (date, time)
timezone = (minute, hour)
```
