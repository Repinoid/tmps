package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

var delays = []int{1, 3, 5}

// func main1() {
func main() {
	postal()
	//	fmt.Printf("response - %+v\nerror %v\n", resp.Header(), err)
}

func postal() {

	httpc := resty.New() //
	httpc.SetBaseURL("http://localhost:8080")

	httpc.SetRetryCount(len(delays))
	delays = delays[1:]
	httpc.SetRetryWaitTime(1 * time.Second)    // начальное время повтора
	httpc.SetRetryMaxWaitTime(9 * time.Second) // 1+3+5
	tn := time.Now()                           // -------------
	httpc.SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		rwt := client.RetryWaitTime
		fmt.Printf("waittime \t%+v\t time %+v  count %d\n", rwt, time.Since(tn), client.RetryCount) // -------
		client.SetRetryWaitTime(time.Duration(delays[0]) * time.Second)
		if len(delays) > 1 {
			delays = delays[1:]
		}
		//	client.SetRetryWaitTime(rwt + 2*time.Second)
		tn = time.Now() // ----------------
		return client.RetryWaitTime, nil
	})
	req := httpc.R().
		SetBody("12345").
		SetHeader("Accept", "text/html").
		SetHeader("Content-Type", "text/html").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("X-Real-IP", "127.0.0.1")

	req.Header.Add("hzz", "WTF")
	resp, err := req.
		SetDoNotParseResponse(false).
		Post("/g")

	log.Printf("StatusCode %d\tErr %v\tResult %s\n", resp.StatusCode(), err, resp)
}
