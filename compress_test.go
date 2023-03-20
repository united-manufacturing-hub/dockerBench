package main

import (
	"bytes"
	_ "embed"
	"github.com/klauspost/compress"
	"github.com/klauspost/compress/zstd"
	"testing"
)

var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil)

// zstdCompress compresses data using zstd compression.
// If the data is not compressible, it will be returned as is with a leading 0.
// If the data is compressible, it will be returned with a leading 1.
func zstdCompress(data []byte) (compressedData []byte) {
	if compress.Estimate(data) <= 0.1 {
		return append([]byte{0}, data...)
	}
	compressed := encoder.EncodeAll(data, make([]byte, 0, len(data)))
	if len(compressed) >= len(data) {
		return append([]byte{0}, data...)
	} else {
		return append([]byte{1}, compressed...)
	}
}

// zstdDecompress decompresses data using zstd compression.
// It strips the leading 0 or 1 and decompresses the data if necessary.
func zstdDecompress(compressedData []byte) (data []byte, err error) {
	if compressedData[0] == 0 {
		return compressedData[1:], nil
	} else {
		return decoder.DecodeAll(compressedData[1:], make([]byte, 0, len(compressedData)))
	}
}

const TestGoodCompression = "The Gopher mascot was introduced in 2009 for the open source launch of the language. The design, by Renée_French, borrowed from a c. 2000 WFMU promotion. [31]\n\nIn November 2016, the Go and Go Mono fonts were released by type designers Charles Bigelow and Kris Holmes specifically for use by the Go project. Go is a humanist sans-serif resembling Lucida Grande, and Go Mono is monospaced. Both fonts adhere to the WGL4 character set and were designed to be legible with a large x-height and distinct letterforms. Both Go and Go Mono adhere to the DIN 1450 standard by having a slashed zero, lowercase l with a tail, and an uppercase I with serifs.[32][33]\n\nIn April 2018, the original logo was replaced with a stylized GO slanting right with trailing streamlines. (The Gopher mascot remained the same.[34])"

//go:embed logs.log
var kafkaLog []byte

func BenchmarkCompressGood(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compressedData := zstdCompress([]byte(TestGoodCompression))
		uncompressedData, err := zstdDecompress(compressedData)
		if err != nil {
			b.Fatal(err)
		}
		if string(uncompressedData) != TestGoodCompression {
			b.Fatal("not equal")
		}
	}
}

func BenchmarkCompressLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compressedData := zstdCompress(kafkaLog)
		uncompressedData, err := zstdDecompress(compressedData)
		if err != nil {
			b.Fatal(err)
		}
		if !bytes.Equal(uncompressedData, kafkaLog) {
			b.Fatal("not equal")
		}
	}
}
