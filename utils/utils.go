package utils

import (
	"math/rand"
	"time"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 可用字元
	textlen  = uint64(len(alphabet))                                            // 可用字元數
	maxlen   = 10                                                               // 最長短網址長度
	minlen   = 3                                                                // 最短短網址長度
	extra    = 2                                                                // textlen/(textlen+extra)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandID() string {
	var id string
	for i := 0; i < minlen; i++ {
		ch := rand.Intn(int(textlen))
		id += string(alphabet[ch])
	}
	for i := minlen; i < maxlen; i++ {
		ch := rand.Intn(int(textlen) + extra)
		if ch >= int(textlen) {
			break
		}
		id += string(alphabet[ch])
	}
	return id
}
