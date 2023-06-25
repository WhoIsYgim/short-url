package generator

import (
	"math/rand"
)

const (
	DefaultAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	DefaultLength   int    = 20
)

type Generator struct {
	hd       GenData
	alphabet []rune
	deleted  chan string
}

type GenData struct {
	Alphabet string
	Length   int
}

func (h *Generator) GenString() string {
	outRunes := make([]rune, 0, h.hd.Length)
	for i := 0; i < h.hd.Length; i++ {
		// do not provide special source because rand.Rand is not goroutine-safe
		idx := rand.Int() % len(h.alphabet)
		outRunes = append(outRunes, h.alphabet[idx])
	}
	return string(outRunes)
}

func NewGenerator() *Generator {
	return &Generator{
		hd: GenData{
			DefaultAlphabet,
			DefaultLength,
		},
		alphabet: []rune(DefaultAlphabet),
	}
}

func NewGeneratorWithData(hd GenData) *Generator {
	return &Generator{
		hd: GenData{
			Alphabet: hd.Alphabet,
			Length:   hd.Length,
		},
		alphabet: []rune(hd.Alphabet),
	}
}
