package main

import (
	"context"
	"errors"
	"fmt"
	"internal/dbaser"
	"log"
	"os"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func main() {
	ctx := context.Background()
	url := os.Getenv("DATABASE_DSN")

	url, _ = os.LookupEnv("DATABASE_DSN")

	//url = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
	url = "postgres://postgres:passwordas@localhost:5432/forgot"
	//	postgres://postgres:mypassword@rds-postgres.xxxxx.amazonaws.com:5432
	//	postgres://postgres:zalupa77@rds-postgres.xxxxx.amazonaws.com:5432

	//db, err := sql.Open("pgx", url)

	db, err := pgx.Connect(ctx, url)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("err code is %v", pgErr.Code)
			pgErr.
		}

		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	mname := "aaaa"
	tx, err := db.Begin(ctx)
	if err != nil {
		fmt.Printf("no db.Begin   %v\n", err)
	}
	order := fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]f);", mname, 3456.77)
	_, err = tx.Exec(ctx, order)
	if err != nil {
		// если ошибка, то откатываем изменения
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("metric %s already exists", mname)
		}

		fmt.Printf("bad tx.exec %s %v\n", mname, err)
		tx.Rollback(ctx)
	}
	tx.Commit(ctx)

	f, err := dbaser.TableGetGauge(ctx, db, mname)
	if err != nil {
		fmt.Printf("bad get %s %v\n", mname, err)
	} else {
		fmt.Printf("value of %s is %f\n", mname, f)
	}

	//perr := db.Ping(ctx)
	//log.Printf("perr %v\n conn port %v\n host %v\n", perr, db.Config().DefaultQueryExecMode, db.Config().Host)

	//crea := "CREATE TABLE IF NOT EXISTS t1 (c1 INT, c2 VARCHAR(10) ); "
	// err = dbaser.TableCreation(ctx, db)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to create tables: %v\n", err)
	// 	os.Exit(1)
	// }
	// err = dbaser.TablePutGauge(ctx, db, "Alloc", 33.88)
	// if err != nil {
	// 	log.Printf("update err %v\n", err)
	// }
	// err = dbaser.TablePutCounter(ctx, db, "Someint", 188)
	// if err != nil {
	// 	log.Printf("update err %v\n", err)
	// }
	// mname := "Alloc"
	// f, err := dbaser.TableGetGauge(ctx, db, mname)
	// if err != nil {
	// 	log.Printf("bad get %s %v\n", mname, err)
	// } else {
	// 	fmt.Printf("value of %s is %f\n", mname, f)
	// }
	// mname = "Someint"
	// i, err := dbaser.TableGetCounter(ctx, db, mname)
	// if err != nil {
	// 	log.Printf("bad get %s %v\n", mname, err)
	// } else {
	// 	fmt.Printf("value of %s is %d\n", mname, i)
	// }
	// str := "SELECT tablename, schemaname FROM pg_catalog.pg_tables WHERE schemaname = $1;"
	// //result, err := db.QueryRow   (ctx, str)
	// rows, err := db.Query(ctx, str, "public")
	// if err != nil {
	// 	log.Printf("bad query %v\n", err)
	// }

	// var tablename string
	// var tablespace string
	// for rows.Next() {
	// 	err = rows.Scan(&tablename, &tablespace)
	// 	if err != nil {
	// 		log.Printf("bad scan\n %v\n", err)
	// 	}
	// 	fmt.Printf("%s %s \n", tablename, tablespace)
	// }
}
