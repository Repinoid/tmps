package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	//	"github.com/jackc/pgx"
)

func main() {
	ctx := context.Background()
	url := os.Getenv("DATABASE_DSN")

	url, _ = os.LookupEnv("DATABASE_DSN")

	url = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
	//url = "postgres://postgres:passwordas@localhost:5432/forgo"
	//	postgres://postgres:mypassword@rds-postgres.xxxxx.amazonaws.com:5432
	//	postgres://postgres:zalupa77@rds-postgres.xxxxx.amazonaws.com:5432

	//db, err := sql.Open("pgx", url)

	db, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	//perr := db.Ping(ctx)
	//log.Printf("perr %v\n conn port %v\n host %v\n", perr, db.Config().DefaultQueryExecMode, db.Config().Host)

	//crea := "CREATE TABLE IF NOT EXISTS t1 (c1 INT, c2 VARCHAR(10) ); "
	err = tableCreation(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create tables: %v\n", err)
		os.Exit(1)
	}
	err = tablePutGauge(ctx, db, "Alloc", 33.88)
	if err != nil {
		log.Printf("update err %v\n", err)
	}
	err = tablePutCounter(ctx, db, "Someint", 188)
	if err != nil {
		log.Printf("update err %v\n", err)
	}
	mname := "Alloc"
	f, err := tableGetGauge(ctx, db, mname)
	if err != nil {
		log.Printf("bad get %s %v\n", mname, err)
	} else {
		fmt.Printf("value of %s is %f\n", mname, f)
	}
	mname = "Someint"
	i, err := tableGetCounter(ctx, db, mname)
	if err != nil {
		log.Printf("bad get %s %v\n", mname, err)
	} else {
		fmt.Printf("value of %s is %d\n", mname, i)
	}
	str := "SELECT tablename, schemaname FROM pg_catalog.pg_tables WHERE schemaname = $1;"
	//result, err := db.QueryRow   (ctx, str)
	rows, err := db.Query(ctx, str, "public")
	if err != nil {
		log.Printf("bad query %v\n", err)
	}

	var tablename string
	var tablespace string
	for rows.Next() {
		err = rows.Scan(&tablename, &tablespace)
		if err != nil {
			log.Printf("bad scan\n %v\n", err)
		}
		fmt.Printf("%s %s \n", tablename, tablespace)
	}
}
