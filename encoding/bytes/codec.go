package bytes

import "github.com/rprtr258/sixty"

// Codec round-trips schema values through the byte-aligned format.
var Codec = sixty.CodecFunc(encode, decode)
