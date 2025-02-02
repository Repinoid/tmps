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

var dbEndPoint = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
var UsersTable = "accounts"
var OrdersTable = "orders"
var TokensTable = "tokens"

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

func (dataBase *DBstruct) UsersTableCreation(ctx context.Context) error {
	db := dataBase.DB
	// В PostgreSQL нельзя передавать название таблицы в качестве параметра
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + UsersTable +
			"(id INT GENERATED ALWAYS AS IDENTITY UNIQUE," +
			"login VARCHAR(100) PRIMARY KEY," +
			"password VARCHAR(200)," +
			"user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"

	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create users table. %w", err)
	}
	return nil
}
func (dataBase *DBstruct) OrdersTableCreation(ctx context.Context) error {
	db := dataBase.DB
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + OrdersTable +
			"(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY," +
			"userCode INT," +
			"orderNumber BIGINT," +
			"order_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
			"FOREIGN KEY (userCode) REFERENCES " + UsersTable + "(id) ON DELETE CASCADE);"

	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create orders table. %w", err)
	}
	return nil
}
func (dataBase *DBstruct) TokensTableCreation(ctx context.Context) error {
	db := dataBase.DB
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + TokensTable +
			"(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY," +
			"userCode INT," +
			"balance BIGINT," +
			"bonus BIGINT," +
			"token VARCHAR(1000)," +
			"token_valid_until TIMESTAMP," +
			"token_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
			"FOREIGN KEY (userCode) REFERENCES " + UsersTable + "(id) ON DELETE CASCADE);"
	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create orders table. %w", err)
	}
	return nil
}

func ConnectToDB(ctx context.Context) (*DBstruct, error) {
	DB, err := ConnectUsersTable(ctx, dbEndPoint)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return nil, err
	}
	err = DB.UsersTableCreation(ctx)
	if err != nil {
		fmt.Printf("tbl: %v", err)
		return nil, err
	}
	err = DB.OrdersTableCreation(ctx)
	if err != nil {
		fmt.Printf("tbl: %v", err)
		return nil, err
	}
	err = DB.TokensTableCreation(ctx)
	if err != nil {
		fmt.Printf("tbl: %v", err)
		return nil, err
	}
	return DB, nil
}

func (dataBase *DBstruct) AddUser(ctx context.Context, userName string, password string) error {
	db := dataBase.DB
	order := "INSERT INTO " + UsersTable + " (login, password) VALUES ($1, crypt($2, gen_salt('md5'))) ;"
	_, err := db.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("add user error is %w", err)
	}
	return nil
}
func (dataBase *DBstruct) CheckUserPassword(ctx context.Context, userName, password string) error {
	db := dataBase.DB
	order := "SELECT (password = crypt($2, password)) AS password_match FROM " + UsersTable + " WHERE login= $1 ;"
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
func (dataBase *DBstruct) IfUserExists(ctx context.Context, userName string) error {
	db := dataBase.DB
	order := "SELECT 7 from " + UsersTable + " WHERE login= $1 ;"
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

func (dataBase *DBstruct) ChangePassword(ctx context.Context, userName string, password string) error {
	db := dataBase.DB
	order := "UPDATE " + UsersTable + " SET password = crypt($2, gen_salt('md5')) WHERE login= $1 ;"
	_, err := db.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("change password error %w", err)
	}
	return nil
}

func (dataBase *DBstruct) AddOrder(ctx context.Context, userName string, orderNumber int64) error {
	db := dataBase.DB

	order := fmt.Sprintf("INSERT INTO %s(userCode, ordernumber) VALUES ((select id from %s where login = '%s'), %d) ;", OrdersTable, UsersTable, userName, orderNumber)
	//INSERT INTO Order (userCode, ordernumber) VALUES ((select id from usera where login = 'user2'), 12345) ;
	_, err := db.Exec(ctx, order)
	if err != nil {
		return fmt.Errorf("add ORDER %w", err)
	}
	return nil
}
