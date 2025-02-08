package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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

	// cmnd := exec.Command("./acc.exe", "-d=postgres://postgres:passwordas@localhost:5432/forgo")
	// //cmnd.Run() // and wait
	// cmnd.Start()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	marks := []buyback{{Match: "Bork", Reward: 10, Reward_type: "%"},
		{Match: "Acer", Reward: 20, Reward_type: "pt"},
		{Match: "HP", Reward: 15, Reward_type: "%"},
		{Match: "Samsung", Reward: 25, Reward_type: "%"},
		{Match: "Apple", Reward: 35, Reward_type: "%"},
	}
	for _, r := range marks {
		buyM, _ := json.Marshal(r)
		poster("/api/goods", buyM)
	}

	ord := orda{Order: strconv.Itoa(Luhner(10000)), Goods: []tovar{
		{Description: "Tea Bork", Price: 111},
		{Description: "Monitor Samsung", Price: 2222},
		{Description: "Apple pc", Price: 111},
		{Description: "HP printer", Price: 2222},
	}}
	ord1 := orda{Order: strconv.Itoa(Luhner(30000)), Goods: []tovar{
		{Description: "Tea Bork", Price: 111},
		{Description: "Monitor Acer", Price: 2222},
		{Description: "Apple watch", Price: 111},
		{Description: "Display Samsung", Price: 2222},
	}}
	buyM, _ := json.Marshal(ord)
	poster("/api/orders", buyM)
	buyM, _ = json.Marshal(ord1)
	poster("/api/orders", buyM)

	getorder(10000)
	getorder(20000)
	getorder(30000)

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
	log.Printf("AGENT responce from server %+v  body is %s\n", resp.StatusCode(), resp.Body())
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
