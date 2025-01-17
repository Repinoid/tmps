package main

import (
	"context"
	"fmt"
	"oppa/internal/dbaser"
	"time"

	"github.com/jackc/pgx"
)

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
