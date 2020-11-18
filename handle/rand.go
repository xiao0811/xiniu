package handle

import (
	"math/rand"
	"time"
)

// RandInt 返回一个在 min, max之间的随机整数
func RandInt(min, max int) int {
	if max <= min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandStr 返回一个长度为n的字符串
func RandStr(n int) string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
