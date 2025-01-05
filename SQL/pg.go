package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	//	ctx := context.Background()
	url := os.Getenv("DATABASE_DSN")
	
	url, _ = os.LookupEnv("DATABASE_DSN")
	
	url = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/postgres"
	//url = "postgres://postgres:passwordas@localhost:5432/postgres"
	//	postgres://postgres:mypassword@rds-postgres.xxxxx.amazonaws.com:5432
	//	postgres://postgres:zalupa77@rds-postgres.xxxxx.amazonaws.com:5432

	db, err := sql.Open("pgx", url)

	//	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	perr := db.Ping()
	log.Printf("perr %v\n conn type %T\n value %v\n", perr, db, db)

	// var name string
	// var weight int64
	// err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	// if err != nil {
	// 	//		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(name, weight)
}
