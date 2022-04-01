package encdec

import (
	"errors"
	"strings"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 可用字元
	textlen  = uint64(len(alphabet))                                            // 可用字元數
	maxlen   = 10                                                               // 最長短網址長度
	minlen   = 3                                                                // 最短短網址長度
)

func Encode(number uint64) string { // 轉成62進位字串
	var result string
	for ; number > 0; number = number / textlen {
		result = string(alphabet[(number%textlen)]) + result
	}
	return result
}

func Decode(encoded string) (uint64, error) { // 62位元轉10進位整數
	var number uint64

	for _, symbol := range encoded {
		pos := strings.IndexRune(alphabet, symbol)
		if pos == -1 { // 找不到字元回傳錯誤
			return uint64(pos), errors.New("unknown char: " + string(symbol))
		}
		number = number*textlen + uint64(pos) // 62位元轉10進位整數
	}
	return number, nil // 回傳10進位
}
