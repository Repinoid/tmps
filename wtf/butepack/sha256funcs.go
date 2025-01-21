package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func Ptr[PP int64 | float64](w PP) *PP {
	i := w
	return &i
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
type Gauge float64
type Counter int64


func encryptS2B(stringToEncrypt string, keyString string) (encryptedString []byte, err error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, _ := RandBytes(aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	//	return fmt.Sprintf("%x", ciphertext), nil
	return ciphertext, nil
}

func decryptS2S(encryptedString, keyString string) (decryptedString string, err error) {

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