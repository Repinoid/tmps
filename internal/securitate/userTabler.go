package securitate

import (
	"context"
	"fmt"

	//	"github.com/jackc/pgx/v4"
	pgx "github.com/jackc/pgx/v5"
)

type DBstruct struct {
	DB *pgx.Conn
}

// соединение с базой данных
func ConnectUsersTable(ctx context.Context, dbEndPoint string) (*DBstruct, error) {
	dataBase := &DBstruct{}
	baza, err := pgx.Connect(ctx, dbEndPoint)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB %s err %w", dbEndPoint, err)
	}
	dataBase.DB = baza
	return dataBase, nil
}

func (dataBase *DBstruct) UsersTableCreation(ctx context.Context, tableName string) error { //  task_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	db := dataBase.DB
	// В PostgreSQL нельзя передавать название таблицы в качестве параметра, so Sprintf
	creatorOrder := "CREATE TABLE IF NOT EXISTS " + tableName + " (id INT GENERATED ALWAYS AS IDENTITY,"
	creatorOrder += "login VARCHAR(100) PRIMARY KEY, password VARCHAR(100)) ;"
	tag, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("error create users table. Tag is \"%s\" error is %w", tag.String(), err)
	}
	return nil
}

func (dataBase *DBstruct) AddUser(ctx context.Context, tableName string, userName string, password string) error {
	db := dataBase.DB
	order := "INSERT INTO " + tableName + " (login, password) VALUES ($1, crypt($2, gen_salt('md5'))) ;"
	_, err := db.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("add user error is %w", err)
	}
	return nil
}
func (dataBase *DBstruct) CheckUserPassword(ctx context.Context, tableName, userName, password string) error {
	db := dataBase.DB
	order := "SELECT (password = crypt($2, password)) AS password_match FROM " + tableName + " WHERE login= $1 ;"
	row := db.QueryRow(ctx, order, userName, password) // password here - what was entered
	var yes bool
	err := row.Scan(&yes)
	if err != nil {
		return fmt.Errorf("QueryRow, error is %w", err)
	}
	if !yes {
		return fmt.Errorf("password not match")
	}
	return nil
}

// nil - user exists
func (dataBase *DBstruct) IfUserExists(ctx context.Context, tableName, userName string) error {
	db := dataBase.DB
	order := "SELECT 7 from " + tableName + " WHERE login= $1 ;"
	row := db.QueryRow(ctx, order, userName) // password here - what was entered
	var yes int
	err := row.Scan(&yes)
	if err != nil {
		return fmt.Errorf(" QueryRow, error is %w", err)
	}
	if yes != 7 {
		return fmt.Errorf("user %s does not exist", userName)
	}
	return nil
}

func (dataBase *DBstruct) ChangePassword(ctx context.Context, tableName, userName string, password string) error {
	db := dataBase.DB
	order := "UPDATE " + tableName + " SET password = crypt($2, gen_salt('md5')) WHERE login= $1 ;"
	_, err := db.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("change password error %w", err)
	}
	return nil
}
