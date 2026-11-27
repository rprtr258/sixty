package bit

import "github.com/rprtr258/sixty"

// Codec round-trips schema values through the bit-packed byte format.
var Codec = sixty.CodecFunc(encode, decode)
