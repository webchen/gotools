package randtool

import (
	"math/rand"
	"time"
)

// RandString 随机字符串
func RandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(25) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
