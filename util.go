package main

import "math/rand"

var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

const letterBytes = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 16                   // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
func int16ToByte(src []int16) []byte {
	b := make([]byte, 0)
	for _, i16 := range src {
		var h, l = uint8(i16 >> 8), uint8(i16 & 0xff)
		b = append(b, l)
		b = append(b, h)
	}
	return b
}
