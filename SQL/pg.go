package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"oppa/internal/dbaser"

	"github.com/jackc/pgx/v5"
)

type MemStorage struct {
	Gaugemetr map[string]gauge
	Countmetr map[string]counter
	Mutter    sync.RWMutex
}

var AttemptDelays = []int{1, 3, 5}

type gauge = dbaser.Gauge
type counter = dbaser.Counter
type Metrics = dbaser.Metrics

func main() {
	ctx := context.Background()

	//url = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"

	url := "postgres://postgres:passwordas@localhost:5432/forgo"

	db, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	var intGag int64 = 6
	var floatGag float64 = 345.345

	metrga := Metrics{ID: "aname3", MType: "gauge", Value: &floatGag}
	metrco := Metrics{ID: "aname1", MType: "counter", Delta: &intGag}

	me := []Metrics{metrga, metrco}
	err = dbaser.TableBuncher(ctx, db, &me)
	if err != nil {
		log.Printf("bad bunch\n %v\n", err)
	}

	flo := 6.5
	meme := dbaser.Metrics{ID: "hz", MType: "gauge", Delta: nil, Value: &flo}
	err = dbaser.TableUpSert(ctx, db, &meme)
	if err != nil {
		log.Printf("bad ONSERT\n %v\n", err)
	}
	meme = dbaser.Metrics{ID: "hzz", MType: "counter", Delta: &intGag, Value: nil}
	err = dbaser.TableUpSert(ctx, db, &meme)
	if err != nil {
		log.Printf("bad ONSERT\n %v\n", err)
	}

	m := []dbaser.Metrics{}
	err = dbaser.TableGetAllTables(ctx, db, &m)
	if err != nil {
		log.Printf("bad allgauges\n %v\n", err)
	}
	fmt.Println(len(m))

	alloc := Metrics{ID: "Alloc", MType: "gauge"}
	err = dbaser.TableGetMetric(ctx, db, &alloc)
	if err != nil {
		log.Printf("bad GET alloc %v\n", err)
	}
	fmt.Printf("alloc %+v  value %f\n", alloc, *alloc.Value)
	poll := Metrics{ID: "PollCount", MType: "counter"}
	err = dbaser.TableGetMetric(ctx, db, &poll)
	if err != nil {
		log.Printf("bad GET alloc %v\n", err)
	}
	fmt.Printf("alloc %+v  value %d\n", poll, *poll.Delta)
	
	poll = Metrics{ID: "aname1", MType: "counter"}
	err = dbaser.TableGetMetric(ctx, db, &poll)
	if err != nil {
		log.Printf("bad GET\n %v\n", err)
	}
	fmt.Println(*poll.Delta)
	
	poll = Metrics{ID: "aname3", MType: "gauge"}
	err = dbaser.TableGetMetric(ctx, db, &poll)
	if err != nil {
		log.Printf("bad GET\n %v\n", err)
	}
	fmt.Println(*poll.Value)
}
