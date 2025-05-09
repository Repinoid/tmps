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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// var isCoded = false
var isCoded = true

func loadClientTLSCredentials(cert string) (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile(cert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	config := &tls.Config{
		// Set InsecureSkipVerify to skip the default validation we are
		// replacing. This will not disable VerifyConnection.
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}
	return credentials.NewTLS(config), nil
}

func main() {

	var conn *grpc.ClientConn
	var err error
	// устанавливаем соединение с сервером
	if isCoded {
		tlsCreds, err := loadClientTLSCredentials("../pems/cert.pem")
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}
		conn, err = grpc.NewClient(":3200", grpc.WithTransportCredentials(tlsCreds))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
	} else {
		// insecure.NewCredentials() - без шифровки
		conn, err = grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
	}

	client := pb.NewMetricClient(conn)

	m := []*pb.Metr{
		{ID: "dd", MType: "counter", Delta: 67},
		{ID: "dd11", MType: "counter", Delta: 67222}}

	md := metadata.New(map[string]string{"token": "12345"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.AddBunch(ctx, &pb.Bunch{
		Meters: m,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.Error != "" {
		fmt.Println(resp.Error)
	}
	fmt.Printf("Client %s\n", resp.OutData)
}
