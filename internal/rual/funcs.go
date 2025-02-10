package rual

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/theplant/luhn"
)

type Tovar struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
type Buyback struct {
	Match       string `json:"match"`
	Reward      int    `json:"reward"`
	Reward_type string `json:"reward_type"`
}
type orda struct {
	Order string  `json:"order"`
	Goods []Tovar `json:"goods"`
}
type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

var accrualhost = "localhost:8080"

// func main() {

// 	cmnd := exec.Command("./acc.exe", "-d=postgres://postgres:passwordas@localhost:5432/forgo")
// 	cmnd.Start()

// 	time.Sleep(time.Second)

// 	if err := run(); err != nil {
// 		panic(err)
// 	}

// }
var marks = []Buyback{
	{Match: "Acer", Reward: 20, Reward_type: "pt"},
	{Match: "Bork", Reward: 10, Reward_type: "%"},
	{Match: "Asus", Reward: 20, Reward_type: "pt"},
	{Match: "Samsung", Reward: 25, Reward_type: "%"},
	{Match: "Apple", Reward: 35, Reward_type: "%"},
}

func LoadGood(num int, goodIdx int, price float64) error {
	ord := orda{Order: strconv.Itoa(Luhner(num)), Goods: []Tovar{
		{Description: "Smth " + marks[goodIdx].Match + " " + strconv.Itoa(num), Price: price}}}
	buyM, _ := json.Marshal(ord)
	err := poster("/api/orders", buyM)
	return err
}

func poster(postCMD string, wts []byte) error {
	httpc := resty.New() //
	httpc.SetBaseURL("http://" + accrualhost)
	req := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(wts)
	resp, err := req.
		SetDoNotParseResponse(false).
		Post(postCMD) //
	log.Printf("%s responce from server %+v  body is %s\n", postCMD, resp.StatusCode(), resp.Body())
	return err
}

func Luhner(numb int) int {
	// if luhn.Valid(numb) {
	// 	return numb
	// }
	return 10*numb + luhn.CalculateLuhn(numb)
}

// OrderStatus - {номер заказа; статус расчёта начисления; рассчитанные баллы к начислению}
func GetFromAccrual(number string) (orderStat OrderStatus, StatusCode int, err error) {
	httpc := resty.New() //
	httpc.SetBaseURL("http://" + accrualhost)
	getReq := httpc.R()

	resp, err := getReq.
		SetResult(&orderStat).
		SetDoNotParseResponse(false).
		SetHeader("Content-Type", "application/json").
		Get("/api/orders/" + number)

	return orderStat, resp.StatusCode(), err
}
