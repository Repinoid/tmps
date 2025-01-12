package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"internal/dbaser"

	"github.com/jackc/pgx/v5"
	//	"github.com/jackc/pgx"
)

var AttemptDelays = []int{1, 3, 5}

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
	m := map[string]float64{}
	err = TableWrapper[float64](dbaser.TableGetAllTables)(ctx, db, &m)
	//err = fu(ctx, db, &m)
	//	err = TableWrapper(dbaser.TableGetAllCounters(ctx, db, &m))
	//err = dbaser.TableGetAllCounters(ctx, db, &m)
	if err != nil {
		log.Printf("bad allgauges\n %v\n", err)
	}
	fmt.Println(len(m))
	mi := map[string]int64{}
	err = TableWrapper[int64](dbaser.TableGetAllTables)(ctx, db, &mi)
	if err != nil {
		log.Printf("bad allgauges\n %v\n", err)
	}
	fmt.Println("countr", len(mi))
}

//func TableGetAllCounters[T Number](ctx context.Context, db *pgx.Conn, mappa *map[string]T) error

func TableWrapper[MV dbaser.MetricValueTypes](origFunc func(ctx context.Context, db *pgx.Conn, mappa *(map[string]MV)) error) func(ctx context.Context,
	db *pgx.Conn, mappa *(map[string]MV)) error {
	wrappedFunc := func(ctx context.Context, db *pgx.Conn, mappa *(map[string]MV)) error {

		err := origFunc(ctx, db, mappa)
		if err != nil {
			for _, delay := range AttemptDelays {
				time.Sleep(time.Duration(delay) * time.Second)
				if err = origFunc(ctx, db, mappa); err == nil {
					break
				}
				fmt.Println(delay, " wrapped !")
			}
		}
		return err
	}
	return wrappedFunc

}
