package common

import (
	"compress/gzip"
	"math/rand"
	"sync"
)

var GzipWriterPool *sync.Pool

func init() {
	GzipWriterPool = &sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
}

var letters = []byte("abcdefghijkmnpqrstuvwxyz") // o,l excluded
var consonants = []byte("bcdfghjkmnpqrstvwxyz")  // l excluded
var vowels = []byte("aeiu")                      // o excluded
var numbers = []byte("23456789")                 // 0,1 excluded

// AccessCode: create access code (ab9cd) code, ex: "bx3ck"
func AccessCode() string {
	code := make([]byte, 7)
	for i := 0; i < 7; i++ {
		switch i {
		case 0, 2, 4, 6:
			code[i] = consonants[rand.Intn(len(consonants))]
		case 1, 5:
			code[i] = vowels[rand.Intn(len(vowels))]
		case 3:
			code[i] = numbers[rand.Intn(len(numbers))]
		}
	}
	return string(code)
}

// RandCode: create random code, with 1 position a number and the rest letters
func RandCode(length int) string {
	numPosition := rand.Intn(length) // random position of number (0 to length-1)
	for numPosition == 0 {
		numPosition = rand.Intn(length) // num postion cannot be zero
	}
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		if i == numPosition {
			code[i] = numbers[rand.Intn(len(numbers))]
		} else {
			code[i] = letters[rand.Intn(len(letters))]
		}
	}
	return string(code)
}
