package murmur3

import "github.com/spaolacci/murmur3"

var hasher = murmur3.New64()

// Hash 对传入的参数求murmur3 hash值
func Hash(lurl string) uint64 {
	hasher.Write([]byte(lurl))
	return hasher.Sum64()
}
