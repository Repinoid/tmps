package main

import (

	//	"encoding/base64"

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

	dst := makeHash(nil, bunchOnMarsh, keyB)
	fmt.Printf("\nsha ! %x LEN %d\n", dst, len(dst))

	//your secret textB

	//encryption
	encrypted, _ := encryptB2B(bunchOnMarsh, keyB)
	enhex := fmt.Sprintf("%x", encrypted)
	fmt.Printf("encrypted data: %s\n", enhex)

	bu := []Metrics{}
	//decryption
	//decrypted, _ := decryptS2S(enhex, string(keyB))
	decrypted, _ := decryptB2B(encrypted, keyB)
	fmt.Printf("decrypted data: %s\n", decrypted)

	err := json.Unmarshal(decrypted, &bu)
	fmt.Printf("metr %+v err %v\n", bu, err)
}
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

// func encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {

// generate a random 32 byte key
// func GenerateRandomKey() (string, error) {
// 	b, _ := RandBytes(32)
// 	key := hex.EncodeToString(b)
// 	return key, nil
// }
