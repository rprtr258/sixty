// Package zigzag maps signed integers to unsigned so that small absolute
// values are encoded as small varints.
package zigzag

// Encode maps a signed integer to its unsigned zigzag representation:
// 0, -1, 1, -2, 2, ... become 0, 1, 2, 3, 4, ...
func Encode(x int64) uint64 {
	return uint64((x >> 63) ^ (x << 1))
}

// Decode is the inverse of Encode.
func Decode(x uint64) int64 {
	return int64(x>>1) ^ -int64(x&1)
}
