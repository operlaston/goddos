package main

import (
	"fmt"
	"hash/maphash"
	"io"
	"strings"
)

type AMSHasher struct {
	seed maphash.Seed
}

func NewAMSHasher() *AMSHasher {
	return &AMSHasher{
		seed: maphash.MakeSeed(),
	}
}

func (h *AMSHasher) Hash(item []byte) int {
	var mh maphash.Hash
	mh.SetSeed(h.seed)
	mh.Write(item)

	hashVal := mh.Sum64()

	if hashVal&1 == 1 {
		return 1
	}
	return -1
}

func ams() int {
	r := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasher := NewAMSHasher()

	z := 0
	b := make([]byte, 1)
	for {
		_, err := r.Read(b) // might want to handle n later
		if err == io.EOF {
			break
		}
		fmt.Println(hasher.Hash(b))
		z += hasher.Hash(b)
	}
	fmt.Println(z)

	return z * z
}

func main() {
	fmt.Printf("Estimated f_2 is %v\n", ams())
}
