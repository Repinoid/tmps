package main

import (
	"context"
	"fmt"
)

const dbEndPoint = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"

func main() {

	ctx := context.Background()
	db, err := ConnectUsersTable(ctx, dbEndPoint)

	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}
	err = db.UsersTableCreation(ctx)
	if err != nil {
		fmt.Printf("error  table creation %v", err)
		return
	}
	// err = db.AddUser(ctx, "user2", "pass2")
	// if err != nil {
	// 	fmt.Printf("error user add %v", err)
	// 	return
	// }
	err = db.CheckUserPassword(ctx, "user2", "pass")
	if err != nil {
		fmt.Printf("error user check password %v", err)
		return
	} else {
		fmt.Println("OK password")
	}
	err = db.ChangePassword(ctx, "user2", "pass")
	if err != nil {
		fmt.Printf("error user check password %v", err)
		return
	} else {
		fmt.Println("OK password change")
	}
	err = db.CheckUserPassword(ctx, "user2", "pass")
	if err != nil {
		fmt.Printf("error user check password %v", err)
		return
	} else {
		fmt.Println("NEW OK password")
	}
	err = db.IfUserExists(ctx, "user2")
	if err != nil {
		fmt.Printf("User not exists %v", err)
		return
	} else {
		fmt.Println("User exists")
	}


	
}
