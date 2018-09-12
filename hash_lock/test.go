package main

import (
	"math/rand"
)

// 依据28原则的数据
func genTestData(count int, keyLen int, valueLen int) []KV {
	keyCount := count / 8
	if keyCount < 10 {
		panic("count must greater 80")
	}

	keys := make([][]byte, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = randBytes(keyLen)
	}

	ret := make([]KV, count)
	for i := 0; i < count; i++ {
		ret[i] = KV{
			key:   keys[genIndex(keyCount)],
			value: randBytes(valueLen),
		}
	}
	return ret
}

func genIndex(count int) int {
	count20 := count / 5
	count80 := count - count20

	if rand.Intn(100) < 80 {
		return rand.Intn(count20)
	}
	return count20 + rand.Intn(count80)
}

func randBytes(l int) []byte {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bs := make([]byte, l)
	for i := 0; i < l; i++ {
		bs[i] = str[rand.Intn(len(str))]
	}
	return bs
}
