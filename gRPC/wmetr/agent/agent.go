package main

import (
	// ...

	"context"
	"fmt"
	"log"
	pb "metr/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// устанавливаем соединение с сервером
	//	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewMetricClient(conn)

	m := []*pb.GMetr{{ID: "dd", MType: "counter", Delta: 67}, {ID: "dd11", MType: "counter", Delta: 67222}}

	resp, err := c.AddBunch(context.Background(), &pb.MBunch{
		Bunch: m,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.Error != "" {
		fmt.Println(resp.Error)
	}
	fmt.Printf("Client %s\n", resp.OutData)
}
