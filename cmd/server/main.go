package main

import (
	"context"
	"fmt"
	"net/http"
	"oppa/internal/securitate"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var host = "localhost:8088"
var usersTable = "accounts"

const dbEndPoint = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"

var sugar zap.SugaredLogger
var ctx context.Context
var DB *securitate.DBstruct
var Token string

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var err error
	ctx = context.Background()

	DB, err = securitate.ConnectUsersTable(ctx, dbEndPoint)

	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return err
	}
	err = DB.UsersTableCreation(ctx)
	if err != nil {
		fmt.Printf("error  table creation %v", err)
		return err
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", registerUser).Methods("POST")
	router.HandleFunc("/api/user/login", loginUser).Methods("POST")

	return http.ListenAndServe(host, router)
}

// curl localhost:8088/api/user/register -H "Content-Type":"application/json" -d "{\"login\":\"user1\",\"password\":\"thePass\"}"
// curl localhost:8088/api/user/login -H "Content-Type":"application/json" -d "{\"login\":\"user1\",\"password\":\"thePass\"}"
