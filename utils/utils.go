package utils

import "math/rand"

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 可用字元
	textlen  = uint64(len(alphabet))                                            // 可用字元數
	maxlen   = 10                                                               // 最長短網址長度
	minlen   = 3                                                                // 最短短網址長度
)

func RandID() string {
	n := minlen + rand.Intn(maxlen-minlen+1)
	var id string
	for k := 0; k < n; k++ {
		ch := rand.Intn(int(textlen))
		id += string(alphabet[ch])
	}
	return id
}
