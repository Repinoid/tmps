package main

import (
	// импортируем пакет со сгенерированными protobuf-файлами

	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "metr/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UsersServer поддерживает все необходимые методы сервера.
type MetricServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricServer

	// используем sync.Map для хранения пользователей
	//users sync.Map
}

func main() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricServer(s, &MetricServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

// AddUser реализует интерфейс добавления пользователя.
func (s *MetricServer) AddBunch(ctx context.Context, in *pb.MBunch) (*pb.BunchResponse, error) {
	var response pb.BunchResponse

	b := in.Bunch

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
			fmt.Printf("token %+v\n", values[0])
		}
	}

	return &response, nil
}
