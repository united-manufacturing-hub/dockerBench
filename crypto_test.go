package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
	"sync"
	"testing"
)

func BenchmarkChaCha20Single(b *testing.B) {
	// We don't bench random gen here, but encryption
	b.StopTimer()
	dataToEncrypt := genRandom(OneGigaByte)
	dataToEncryptCopy := make([]byte, OneGigaByte)
	copy(dataToEncryptCopy, dataToEncrypt)

	key := genRandom(chacha20poly1305.KeySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		b.Fatalf("Failed to create ChaCha20Poly1305 cipher: %v", err)
	}

	nonces := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+OneGigaByte+aead.Overhead())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			b.Fatalf("Failed to generate nonce")
		}
		nonces[i] = nonce
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		nonce := nonces[i]
		ciphertext := aead.Seal(nonce, nonce, dataToEncrypt, nil)
		var restoredPlaintext []byte
		restoredPlaintext, err = aead.Open(nil, nonce, ciphertext[aead.NonceSize():], nil)

		if err != nil {
			b.Fatalf("Failed to decrypt: %v", err)
		}
		if len(restoredPlaintext) != OneGigaByte {
			b.Fatalf("Decrypted data length is wrong")
		}
		if !bytes.Equal(restoredPlaintext, dataToEncryptCopy) {
			b.Fatalf("Decrypted data is wrong")
		}
	}
}

func BenchmarkChaCha20Threaded(b *testing.B) {
	// We don't bench random gen here, but encryption
	b.StopTimer()
	dataToEncrypt := genRandom(OneGigaByte)
	dataToEncryptCopy := make([]byte, OneGigaByte)
	copy(dataToEncryptCopy, dataToEncrypt)

	key := genRandom(chacha20poly1305.KeySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		b.Fatalf("Failed to create ChaCha20Poly1305 cipher: %v", err)
	}

	nonces := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+OneGigaByte+aead.Overhead())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			b.Fatalf("Failed to generate nonce")
		}
		nonces[i] = nonce
	}
	b.StartTimer()

	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go encryptDecryptChaCha20Poly1305(aead, nonces[i], dataToEncrypt, &dataToEncryptCopy, &wg)
	}
	wg.Wait()
}

func encryptDecryptChaCha20Poly1305(aead cipher.AEAD, nonce []byte, dataToEncrypt []byte, dataToEncryptCopy *[]byte, s *sync.WaitGroup) {
	defer s.Done()
	ciphertext := aead.Seal(nonce, nonce, dataToEncrypt, nil)
	var restoredPlaintext []byte
	var err error
	restoredPlaintext, err = aead.Open(nil, nonce, ciphertext[aead.NonceSize():], nil)

	if err != nil {
		panic("Failed to decrypt")
	}
	if len(restoredPlaintext) != OneGigaByte {
		panic("Decrypted data length is wrong")
	}
	if !bytes.Equal(restoredPlaintext, *dataToEncryptCopy) {
		panic("Decrypted data is wrong")
	}
}
