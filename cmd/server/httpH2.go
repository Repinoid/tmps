package main

import (
	"log"
	"net/http"
	"oppa/internal/rual"

	"github.com/gorilla/mux"
)

func GetOrders(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	namba := vars["number"]
	var orderStat rual.OrderStatus

	orderStat, statCode, err := rual.GetFromAccrual(namba)

	rwr.WriteHeader(statCode) // return statuscode from accrual

	

	log.Printf("GET %s order %+v  body is %+v err is %+v\n", namba, statCode, orderStat, err)

}
