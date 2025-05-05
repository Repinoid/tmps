package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/g", handleRequest).Methods("POST")

	http.ListenAndServe("localhost:8080", router)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// CIDR in server params
	trusted := "127.0.0.1/24"
	// check CIDR
	_, ipnet, err := net.ParseCIDR(trusted)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	// get X-Real-IP from agent request header
	agentIP, err := ipFromHeader(r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("X-Real-IP %+v\n", agentIP)

	aIP := net.ParseIP(agentIP.String())

	if ipnet.Contains(aIP) {
		log.Printf("IP %s is in CIDR %s", agentIP.String(), trusted)
		w.WriteHeader(200)
		w.Write([]byte(agentIP))
		return

	} else {
		log.Printf("IP %s is NOT in CIDR %s", agentIP.String(), trusted)
		w.WriteHeader(400)
		w.Write([]byte("IP " + agentIP.String() + "not in CIDR" + trusted))
		return
	}

}

func ipFromHeader(r *http.Request) (net.IP, error) {
	// смотрим заголовок запроса X-Real-IP
	ipStr := r.Header.Get("X-Real-IP")
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
}
