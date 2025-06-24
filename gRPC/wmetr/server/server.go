package main

import (
	// импортируем пакет со сгенерированными protobuf-файлами

	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	pb "gorono/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var isCoded = false

//var isCoded = true

// UsersServer поддерживает все необходимые методы сервера.
type MetricServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName> для совместимости с будущими версиями
	pb.UnimplementedMetricServer
}

// loadTLSCredentials загрузка сертификатов
func loadTLSCredentials(cert, key string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	var srv *grpc.Server
	if isCoded {
		// Load TLS credentials
		creds, err := loadTLSCredentials("../pems/cert.pem", "../pems/key.pem")
		if err != nil {
			log.Fatalf("failed to load TLS credentials: %v", err)
		}
		srv = grpc.NewServer(grpc.Creds(creds))
	} else {
		// без шифровки
		srv = grpc.NewServer()
	}
	// регистрируем сервис
	pb.RegisterMetricServer(srv, &MetricServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := srv.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

// 
func (s *MetricServer) AddBunch(ctx context.Context, in *pb.Bunch) (*pb.BunchResponse, error) {
	var response pb.BunchResponse

	b := in.Meters

	bout, err := json.Marshal(b)

	if err != nil {
		return nil, err
	}
	response.OutData = string(bout)

	for i, metr := range b {

		fmt.Printf("server %d %+v\n", i, metr)

	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			fmt.Printf("token %+v len %d\n", values[0], len(values))
		}
		values = md.Get("hz")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			fmt.Printf("hz %+v len %d\n", values[0], len(values))
		}
	}

	return &response, nil
}
