package utils

import "math/rand/v2"

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" //62
	
)

// 生成一个随机字符串
func RandString(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		idx := rand.IntN(len(letters))
		b[i] = letters[idx]
	}
	return string(b)
}
