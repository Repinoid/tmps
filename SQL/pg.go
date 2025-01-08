package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"internal/dbaser"

	"github.com/jackc/pgx/v5"
	//	"github.com/jackc/pgx"
)

func main() {
	ctx := context.Background()
	url := os.Getenv("DATABASE_DSN")

	url, _ = os.LookupEnv("DATABASE_DSN")

	//url = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
	url = "postgres://postgres:passwordas@localhost:5432/forgo"
	//	postgres://postgres:mypassword@rds-postgres.xxxxx.amazonaws.com:5432
	//	postgres://postgres:zalupa77@rds-postgres.xxxxx.amazonaws.com:5432

	//db, err := sql.Open("pgx", url)

	db, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	
//	db.BeginTx()

	m, err := dbaser.TableGetAllGauges(ctx, db)
	if err != nil {
		log.Printf("bad allgauges\n %v\n", err)
	}
	fmt.Println(m)
	mi, err := dbaser.TableGetAllCounters(ctx, db)
	if err != nil {
		log.Printf("bad allgauges\n %v\n", err)
	}
	fmt.Println(mi)

}
