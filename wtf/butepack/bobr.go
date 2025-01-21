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

// encryption
func encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {

	//convert decode it to bytes
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Create a  RANDOM nonce.
nonce, _ := RandBytes(aesGCM.NonceSize())


	// nonce := make([]byte, aesGCM.NonceSize())
	// _, err = io.ReadFull(rand.Reader, nonce)
	// if err != nil {
	// 	return "", err
	// }

	//Encrypt the data
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
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// generate a random 32 byte key
func GenerateRandomKey() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b) // записываем байты в слайс b
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
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