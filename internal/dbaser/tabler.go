package dbaser

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type MetricValueTypes interface {
	int64 | float64
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func TableOnSert(ctx context.Context, db *pgx.Conn, metr Metrics) error {

	order := fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]g) ", metr.ID, *metr.Value)
	order += fmt.Sprintf("ON CONFLICT (metricname) DO UPDATE SET metricname='%[1]s', value=%[2]g;", metr.ID, *metr.Value)
	tag1, err := db.Exec(ctx, order)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error UPDATE Gauge %s with %g value. Tag is \"%s\" error is %w",
		metr.ID, *metr.Value, tag1.String(), err)

}
func TableGetAllTables(ctx context.Context, db *pgx.Conn, metro *([]Metrics)) error {
	zapros := `select 'counter' AS metrictype, metricname AS name, null AS value, value AS delta from counter
		UNION
	select 'gauge' AS metrictype, metricname as name, value as value, null as delta from gauge
	`
	var inta int64
	var flo float64
	metr := Metrics{ID: "", MType: "", Value: &flo, Delta: &inta}

	rows, err := db.Query(ctx, zapros)
	if err != nil {
		return fmt.Errorf("error Query %[2]s:%[3]d  %[1]w", err, db.Config().Host, db.Config().Port)
	}
	for rows.Next() {
		err = rows.Scan(&metr.MType, &metr.ID, &metr.Value, &metr.Delta)
		if err != nil {
			return fmt.Errorf("error table Scan %[2]s:%[3]d  %[1]w", err, db.Config().Host, db.Config().Port)
		}
		(*metro) = append((*metro), metr)
	}
	if err:= rows.Err(); err != nil {
		return err
	}
	return nil
}
//-------------------------------------------------------------------------------------------------------------
func TableGetAllGauges(ctx context.Context, db *pgx.Conn, mappa *(map[string]float64)) error {
	var flo float64
	var str string

	zapros := "SELECT * FROM gauge;"
	rows, err := db.Query(ctx, zapros)
	if err != nil {
		return fmt.Errorf("error Query %[2]s:%[3]d database  %[1]w", err, db.Config().Host, db.Config().Port)
	}
	for rows.Next() {
		err = rows.Scan(&str, &flo)
		if err != nil {
			return fmt.Errorf("error gauge table Scan %[2]s:%[3]d database\n%[1]w", err, db.Config().Host, db.Config().Port)
		}
		(*mappa)[str] = flo
	}
	return nil
}

func TableCreation(ctx context.Context, db *pgx.Conn) error {
	crea := "CREATE TABLE IF NOT EXISTS Gauge(metricname VARCHAR(30) PRIMARY KEY, value FLOAT8);"
	tag, err := db.Exec(ctx, crea)
	if err != nil {
		return fmt.Errorf("error create Gauge table. Tag is \"%s\" error is %w", tag.String(), err)
	}
	crea = "CREATE TABLE IF NOT EXISTS Counter(metricname VARCHAR(30) PRIMARY KEY, value BIGINT);"
	tag, err = db.Exec(ctx, crea)
	if err != nil {
		return fmt.Errorf("error create Counter table. Tag is \"%s\" error is %w", tag.String(), err)
	}
	return nil
}

func TablePutGauge(ctx context.Context, db *pgx.Conn, mname string, value float64) error {

	order := fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]f);", mname, value)
	tag1, err := db.Exec(ctx, order)
	if err == nil {
		return nil
	}
	order = fmt.Sprintf("UPDATE Gauge SET value=%[2]f WHERE metricname='%[1]s'", mname, value)
	tag2, err := db.Exec(ctx, order)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error UPDATE Gauge %s with %f value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %w",
		mname, value, tag1.String(), tag2.String(), err)
}

func TableGetGauge(ctx context.Context, db *pgx.Conn, mname string) (float64, error) {
	var flo float64
	str := "SELECT value FROM gauge WHERE metricname = $1;"
	row := db.QueryRow(ctx, str, mname)
	err := row.Scan(&flo)
	if err != nil {
		return 0.0, fmt.Errorf("error get %s gauge metric.  %w", mname, err)
	}
	return flo, nil
}

func TablePutCounter(ctx context.Context, db *pgx.Conn, mname string, value int64) error {
	order := fmt.Sprintf("INSERT INTO Counter(metricname, value) VALUES ('%[1]s',%[2]d);", mname, value)
	tag1, err := db.Exec(ctx, order)
	if err == nil {
		return nil
	}
	order = fmt.Sprintf("UPDATE Counter SET value=%[2]d WHERE metricname='%[1]s'", mname, value)
	tag2, err := db.Exec(ctx, order)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error UPDATE Counter %s with %d value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %w",
		mname, value, tag1.String(), tag2.String(), err)
}

func TableGetCounter(ctx context.Context, db *pgx.Conn, mname string) (int64, error) {
	var inta int64
	str := "SELECT value FROM counter WHERE metricname = $1;"
	row := db.QueryRow(ctx, str, mname)
	err := row.Scan(&inta)
	if err != nil {
		return 0, fmt.Errorf("error get %s counter metric.  %w", mname, err)
	}
	return inta, nil
}

type Gauge float64
type Counter int64

func TableBunchGauges(ctx context.Context, db *pgx.Conn, gaaga map[string]Gauge) error {
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

func TableBunchCounters(ctx context.Context, db *pgx.Conn, gaaga map[string]Counter) error {
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
