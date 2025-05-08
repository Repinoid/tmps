package main

import (
	// ...

	"context"
	"fmt"
	"log"
	pb "metr/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	md := metadata.New(map[string]string{"token": "12345"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := c.AddBunch(ctx, &pb.MBunch{
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
