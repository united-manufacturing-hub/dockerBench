package main

import "crypto/rand"

func genRandom(size int) []byte {
	randData := make([]byte, size)
	_, _ = rand.Read(randData)
	return randData
}
