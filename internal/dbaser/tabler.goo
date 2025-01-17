package dbaser

import (
	"context"
	"fmt"

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


// -------------------------------------------------------------------------------------------------------------
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

//EXCLUDED в PostgreSQL — это специальная таблица, которая используется для ссылки на значения, предложенные для вставки в команде INSERT.
//Она имеет те же столбцы, что и вставляемая таблица, и её значения — это значения, которые были бы вставлены,
//если бы команда INSERT не столкнулась с конфликтом.
