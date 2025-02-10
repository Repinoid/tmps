package main

import (
	"log"
	"net/http"
	"oppa/internal/rual"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

var accrualhost = "localhost:8080"

func GetOrders(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	namba := vars["number"]
	var orderStat rual.OrderStatus

	httpc := resty.New() //
	httpc.SetBaseURL("http://" + accrualhost)
	getReq := httpc.R()
	// 	SetHeader("Content-Type", "application/json").
	// 	SetBody(wts)
	resp, err := getReq.
		SetResult(&orderStat).
		SetDoNotParseResponse(false).
		SetHeader("Content-Type", "application/json").
		Get("/api/orders/" + namba) 
	rwr.WriteHeader(resp.StatusCode())	// return statuscode from accrual

	


	log.Printf("GET %s order %+v  body is %+v err is %+v\n", namba, resp.StatusCode(), orderStat, err)

}
