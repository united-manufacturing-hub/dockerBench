package main

import (
	"testing"
)

const OneKiloByte = 1024
const OneMegaByte = OneKiloByte * 1024
const OneGigaByte = OneMegaByte * 1024

func BenchmarkGenerateRandomBytes1Kb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genRandom(OneKiloByte)
	}
}
func BenchmarkGenerateRandomBytes1Mb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genRandom(OneMegaByte)
	}
}
func BenchmarkGenerateRandomBytes1Gb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genRandom(OneGigaByte)
	}
}
