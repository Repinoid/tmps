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

type gauge = dbaser.Gauge
type counter = dbaser.Counter

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

	gamap := map[string]gauge{"one": 1.11, "two2": 2.22, "tri3": 3.33}
	err = dbaser.TableBunchGauges(ctx, db, gamap)
	if err != nil {
		fmt.Printf("error ...  %[1]v", err)
	}
	comap := map[string]counter{"one": 1, "two2": 2, "tri3": 3}
	err = dbaser.TableBunchCounters(ctx, db, comap)
	if err != nil {
		fmt.Printf("error ...  %[1]v", err)
	}

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
