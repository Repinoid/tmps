package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main2() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8088`, http.HandlerFunc(webhook))
}

func webhook(rwr http.ResponseWriter, req *http.Request) {
	//	rwr.Header().Set("Content-Type", "text/plain")
	rwr.Header().Set("Content-Type", "application/json")

	outer := fmt.Sprintf("Methos %[1]v Type %[1]T\n", req.Method)
	outer += fmt.Sprintf("URL %[1]v Type %[1]T\n", req.URL.String())
	outer += fmt.Sprintf("Header %[1]v Type %[1]T\n", req.Header)
	outer += fmt.Sprintf("Body %[1]v Type %[1]T\n", req.Body)
	outer += fmt.Sprintf("Host %[1]v Type %[1]T\n", req.Host)
	outer += fmt.Sprintf("ContentLength %[1]v Type %[1]T\n", req.ContentLength)
	outer += fmt.Sprintf("Form %[1]v Type %[1]T\n", req.Form)

	var p []byte
	req.Body.Read(p)
	outer += fmt.Sprintf("BodyRead %[1]v Type %[1]T\n", p)

	urla := req.URL.String()

	splittedURL := strings.Split(urla, "/")
	outer += fmt.Sprintf("Length splitted URL %[1]v Type %[1]T\n", len(splittedURL))

	rwr.WriteHeader(http.StatusOK)
	for i, v := range splittedURL {
		outer += fmt.Sprintf("%d %s\n", i, v)
	}
	rwr.Write([]byte(outer))

}

// curl -v -X POST http://localhost:8088/update/gauge/Alloc/77.88
