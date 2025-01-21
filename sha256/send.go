package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	a:=gener()
	ggg(a)
}

var src = []byte("Видишь гофера? Нет. И я нет. А он есть.")
var key []byte
var err error

func gener() string {
	// подписываемое сообщение

	// создаём случайный ключ
	key, err = generateRandom(16)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ""
	}

	// подписываем алгоритмом HMAC, используя SHA-256
	h := hmac.New(sha256.New, key)
	h.Write(src)
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst)
}

var secretkey = []byte("secret key")
var msg = "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21"

func ggg(msg string) {
	var (
		data []byte // декодированное сообщение с подписью
		id   uint32 // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	id = binary.BigEndian.Uint32(data[:4])
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data[:4])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		fmt.Println("Подпись подлинная. ID:", id)
	} else {
		fmt.Println("Подпись неверна. Где-то ошибка")
	}
}
