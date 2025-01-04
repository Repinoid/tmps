package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func postMetric(metricType, metricName, metricValue string) int {
	host := "http://localhost:8080"
	url := host + "/update/" + metricType + "/" + metricName + "/" + metricValue
	resp, _ := http.Post(url, "text/plain", nil)
	defer resp.Body.Close()
	return resp.StatusCode
}

func main() {
	url := "http://localhost:8088"

	jsonStr := []byte(`{"name":"John","age":30,"city":"New York"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	// Handle response
}

func main2() {
	statusCode := postMetric("gauge", "Alloc", "77.88")
	if statusCode != http.StatusOK {
		log.Fatal(statusCode)
	}
}

func main1() {
	url := "http://localhost:8088/update/gauge/Alloc/77.88"
	fmt.Println("URL:>", url)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	time.Sleep(1 * time.Second)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
