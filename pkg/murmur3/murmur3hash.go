package murmur3

import "github.com/spaolacci/murmur3"

// Hash 对传入的参数求murmur3 hash值
func Hash(lurl string) uint64 {
	return murmur3.Sum64([]byte(lurl))
}
