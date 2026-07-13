package base62

// 62 进制数的表示
// 0-9 A-Z a-z   10 + 26 + 26 = 62

// 10进制转换为 62进制
// 0  -> 0
// 10 -> a
// 35 -> z
// 36 -> A 61 -> Z

// x = x1* 10 ^ 0 + x2 * 10 ^ 1 + x3 * 10 ^ 2 + .... = y1 * 62 ^ 0 + y2 * 62 ^ 1 + y3 * 62 ^ 3
// x % 62 == y1  x = x / 62

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Uint2string 将10进制的seq转换为62进制的字符串表示形式
func Uint2string(seq uint64) string {
	ans := ""

	for seq > 0 {
		ans = string(base62[seq%62]) + ans
		seq /= 62
	}
	return ans
}

// String2Uint 将62进制的字符串转换为对应的10进制数
func String2Uint(s string) uint64 {
	var seq = uint64(0)
	// 正向索引
	for i := 0; i < len(s); i++ {
		seq = seq*62 + uint64(char2value(s[i]))
	}
	return seq
}

// 将字符转换为 0-61 的数值（基于 0-9A-Za-z 顺序）
func char2value(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 10
	}

	if c >= 'a' && c <= 'z' {
		return c - 'a' + 36
	}
	return 0
}
