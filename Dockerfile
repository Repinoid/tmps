FROM golang:alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -C ./cmd/gophermart/  -o ./main.exe .

EXPOSE 8080

CMD ["./cmd/gophermart/main.exe"]
