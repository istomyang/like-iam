package idutil

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestGetRandString(t *testing.T) {
	s, _ := GetRandString(AlphabetL+AlphabetU+Number, 20)
	fmt.Println(s)
}

func TestRandString(t *testing.T) {
	n := 5
	charset := AlphabetL

	var randomness = make([]byte, n)
	read, err := rand.Read(randomness)
	if err != nil || read != n {
		return
	}

	var r = make([]rune, n)
	var csr = []rune(charset)

	for i, rn := range randomness {
		r[i] = csr[int(rn)%len(csr)]
	}

	fmt.Println(r)
	fmt.Println(string(r))
}

func TestRand(t *testing.T) {
	var b = make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return
	}

	// Exp: [147 4 209 128 8 148 241 175 40 169]
	fmt.Println(b)
	fmt.Println(string(b[1]))
}

func TestString(t *testing.T) {
	s := "你好123123jsndajsdfJJSNS"
	fmt.Println([]byte(s))
	fmt.Println([]rune(s))
	fmt.Println(string([]rune(s)[0]))
	fmt.Println(s[0])
}
