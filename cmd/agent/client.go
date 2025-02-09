package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type tovar struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}
type buyback struct {
	Match       string `json:"match"`
	Reward      int    `json:"reward"`
	Reward_type string `json:"reward_type"`
}
type orda struct {
	Order string  `json:"order"`
	Goods []tovar `json:"goods"`
}
type orderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

var host = "localhost:8080"

func main() {

	cmnd := exec.Command("./acc.exe", "-d=postgres://postgres:passwordas@localhost:5432/forgo")
	cmnd.Start()

	time.Sleep(time.Second)

	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {

	marks := []buyback{
		{Match: "Acer", Reward: 20, Reward_type: "pt"},
		{Match: "Bork", Reward: 10, Reward_type: "%"},
		{Match: "Asus", Reward: 20, Reward_type: "pt"},
		{Match: "Samsung", Reward: 25, Reward_type: "%"},
		{Match: "Apple", Reward: 35, Reward_type: "%"},
	}
	for _, r := range marks {
		buyM, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("err %w", err)
		}
		poster("/api/goods", buyM)
	}
// "{\"match\":\"Acer\",\"reward\":20,\"reward_type\":\"pt\"}"
	ordera := []orda{}
	for i := range 30 {
		ord := orda{Order: strconv.Itoa(Luhner(i)), Goods: []tovar{
			{Description: "Smth " + marks[i%4].Match + " " + strconv.Itoa(i), Price: rand.IntN(1000)}, //+ " " + strconv.Itoa(Luhner(i+rand.IntN(777) + 11111))
		}}
		//	log.Printf("desc %s", ord.Order)
		ordera = append(ordera, ord)
	}

	for _, ord := range ordera {
		buyM, _ := json.Marshal(ord)
		poster("/api/orders", buyM)
	}
// "{\"order\":\"0\",\"goods\":[{\"description\":\"Smth Acer 0\",\"price\":729}]}"
	getorder(1)
	getorder(2)
	getorder(3)

	return nil
}

func poster(postCMD string, wts []byte) error {
	httpc := resty.New() //
	httpc.SetBaseURL("http://" + host)
	req := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(wts)
	resp, err := req.
		SetDoNotParseResponse(false).
		Post(postCMD) //
	log.Printf("%s responce from server %+v  body is %s\n", postCMD, resp.StatusCode(), resp.Body())
	return err
}

var orderStat orderStatus

func getorder(number int) error {
	getCMD := fmt.Sprintf("/api/orders/%s", strconv.Itoa(Luhner(number)))
	httpc := resty.New() //
	httpc.SetBaseURL("http://" + host)
	req := httpc.R()
	// 	SetHeader("Content-Type", "application/json").
	// 	SetBody(wts)
	resp, err := req.
		SetResult(&orderStat).
		SetDoNotParseResponse(false).
		Get(getCMD) //
	log.Printf("GET %d order %+v  body is %+v\n", number, resp.StatusCode(), orderStat)
	return err
}
