package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	//	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {

	//generate a random 32 byte key
	key, _ := GenerateRandomKey()

	//your secret text
	secret := "This is my password"

	//encryption
	encrypted, _ := encrypt(secret, key)
	fmt.Printf("encrypted data: %s\n", encrypted)

	//decryption
	decrypted, _ := decrypt((encrypted), (key))
	fmt.Printf("decrypted data: %s\n", decrypted)
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce, _ := RandBytes(aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func decrypt(encryptedString, keyString string) (decryptedString string, err error) {

	key, err := hex.DecodeString(keyString) // hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

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
