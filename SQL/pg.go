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

type gauge float64
type counter int64

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
	err = TableBunchGauges(ctx, db, gamap)
	if err != nil {
		fmt.Printf("error ...  %[1]v", err)
	}
	comap := map[string]counter{"one": 1, "two2": 2, "tri3": 3}
	err = TableBunchCounters(ctx, db, comap)
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

func TableBunchGauges(ctx context.Context, db *pgx.Conn, gaaga map[string]gauge) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error db.Begin  %[1]w", err)
	}
	for gaugeName, value := range gaaga {
		order := fmt.Sprintf("UPDATE Gauge SET value=%[2]f WHERE metricname='%[1]s'", gaugeName, value)
		tagUpdate, _ := tx.Exec(ctx, order)
		tu := tagUpdate.RowsAffected()
		if tu != 0 { // если удалось записать - уже существует и INSERT не нужен
			continue
		}
		order = fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]f);", gaugeName, value)
		tagInsert, err := tx.Exec(ctx, order)
		if err != nil {
			log.Printf("error UPDATE Gauge %s with %f value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %v",
				gaugeName, value, tagInsert.String(), tagUpdate.String(), err)
		}
	}
	return tx.Commit(ctx)
}

func TableBunchCounters(ctx context.Context, db *pgx.Conn, gaaga map[string]counter) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error db.Begin  %[1]w", err)
	}
	for counterName, value := range gaaga {
		order := fmt.Sprintf("UPDATE Counter SET value=%[2]d WHERE metricname='%[1]s'", counterName, value)
		tagUpdate, _ := tx.Exec(ctx, order)
		tu := tagUpdate.RowsAffected()
		if tu != 0 { // если удалось записать - уже существует и INSERT не нужен
			continue
		}
		order = fmt.Sprintf("INSERT INTO Counter(metricname, value) VALUES ('%[1]s',%[2]d);", counterName, value)
		tagInsert, err := tx.Exec(ctx, order)
		if err != nil {
			log.Printf("error UPDATE Counter %s with %d value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %v",
				counterName, value, tagInsert.String(), tagUpdate.String(), err)
		}
	}
	return tx.Commit(ctx)
}
