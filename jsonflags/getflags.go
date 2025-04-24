// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)


type flags struct {
	Address        string `json:"address"`        // аналог переменной окружения ADDRESS или флага -a
	Restore        bool   `json:"restore"`        // аналог переменной окружения RESTORE или флага -r
	Store_interval string `json:"store_interval"` // аналог переменной окружения STORE_INTERVAL или флага -i
	Store_file     string `json:"store_file"`     // аналог переменной окружения STORE_FILE или -f
	Database_dsn   string `json:"database_dsn"`   // аналог переменной окружения DATABASE_DSN или флага -d
	Crypto_key     string `json:"crypto_key"`     // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
	Pusto          string `json:"ahz"`            // pusto
}

func main() {

	content, err := os.ReadFile("flags.json")
	if err != nil {
		log.Fatal(err)
	}

	err = takeParamsFromJSON(content)
	fmt.Println(err)

}
func takeParamsFromJSON(jstring []byte) error {
	var f flags
	err := json.Unmarshal(jstring, &f)
	if err != nil {
		return err
	}

	i, err := takeInterval(f.Store_interval)

	fmt.Printf("err %v %+v\n%d\n", err, f, i)
	return nil
}

// takeInterval возвращает значение интервала в секундах и ошибку
func takeInterval(s string) (t int, err error) {
	if s == "" {
		return 0, nil
	}
	sec, isSec := strings.CutSuffix(s, "s")
	if isSec {
		hm, err := strconv.Atoi(sec)
		if err != nil {
			return 0, err
		}
		return hm, nil
	}
	min, isMin := strings.CutSuffix(s, "m") // на всяк случай минуты
	if isMin {
		hm, err := strconv.Atoi(min)
		if err != nil {
			return 0, err
		}
		return hm * 60, nil
	}
	return 0, errors.New("bad Interval format")
}
