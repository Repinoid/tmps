package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"

	//	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func main() {

	controlMetric := Metrics{MType: "gauge", ID: "Alloc", Value: Ptr[float64](78)}
	//	cmMarshalled, _ := json.Marshal(controlMetric)
	controlMetric1 := Metrics{MType: "gauge", ID: "Alloc", Value: Ptr[float64](77)}
	//	cmMarshalled1, _ := json.Marshal(controlMetric1)

	bunch := []Metrics{controlMetric, controlMetric1}
	bunchOnMarsh, _ := json.Marshal(bunch)

	//key, _ := GenerateByteKey()
	keyB, _ := RandBytes(32)

	ade := []byte("12345")
	h := hmac.New(sha256.New, keyB) // New returns a new HMAC hash using the given hash.Hash type and key.
	h.Write(bunchOnMarsh)           // func (hash.Hash) Sum(b []byte) []byte
	dst := h.Sum(ade)               //Sum appends the current hash to b and returns the resulting slice. It does not change the underlying hash state.
	fmt.Printf("\nsha ! %x LEN %d\n", dst, len(dst))
	fmt.Printf("\nade %s\n", dst[:5])

	//your secret textB

	//encryption
	encrypted, _ := encryptB2B(bunchOnMarsh, keyB)
	enhex := fmt.Sprintf("%x", encrypted)
	fmt.Printf("encrypted data: %s\n", enhex)

	//decryption
	//decrypted, _ := decryptS2S(enhex, string(keyB))
	decrypted, _ := decryptB2B(encrypted, keyB)
	fmt.Printf("decrypted data: %s\n", decrypted)
}

// func encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {

// generate a random 32 byte key
func GenerateRandomKey() (string, error) {
	b, _ := RandBytes(32)
	key := hex.EncodeToString(b)
	return key, nil
}
func GenerateByteKey() (byteKey []byte, err error) {
	rb, err := RandBytes(32)
	byteKey = make([]byte, len(rb)*2)
	n := hex.Encode(byteKey, rb)
	return byteKey[:n], err
}
func RandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
	//  return base64.StdEncoding.EncodeToString(b), nil
}
