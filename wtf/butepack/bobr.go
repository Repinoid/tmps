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


	//secret := "This is my password"

	//generate a random 32 byte key
	key, _ := GenerateRandomKey()

	h := hmac.New(sha256.New, []byte(key))
	h.Write(bunchOnMarsh)
	dst := h.Sum(nil)
	fmt.Printf("\nsha ! %x LEN %d\n", dst, len(dst))

	//your secret text

	//encryption
	encrypted, _ := encryptS2B(string(bunchOnMarsh), key)
	enhex := fmt.Sprintf("%x", encrypted)
	fmt.Printf("encrypted data: %s\n", enhex)

	//decryption
	decrypted, _ := decryptS2S(enhex, (key))
	fmt.Printf("decrypted data: %s\n", decrypted)
}

// func encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {

// generate a random 32 byte key
func GenerateRandomKey() (string, error) {
	b, _ := RandBytes(32)
	key := hex.EncodeToString(b)
	return key, nil
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
