package text

import "github.com/rprtr258/sixty"

// Codec round-trips schema values through the text format.
var Codec = sixty.CodecFunc(encode, decode)
