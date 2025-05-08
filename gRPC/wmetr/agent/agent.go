package main

import (
	// ...

	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	pb "metr/proto"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func loadClientTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile("../pems/cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	// https://pkg.go.dev/crypto/tls#Config
	// Client side configuration.

	creds, err := loadClientTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// устанавливаем соединение с сервером
	//	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(tlsCreds))
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(creds))
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
