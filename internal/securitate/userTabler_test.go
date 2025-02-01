package securitate

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbEndPoint = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
const testTableName = "testable"

var ctx context.Context

func TestDBstruct_AddUser(t *testing.T) {
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name      string
		args      args
		isErr     bool
		errString string
	}{
		{
			name: "Nice adding",
			args: args{
				userName: "us1",
				password: "pass1",
			},
			isErr: false,
		},
		{
			name: "Duplicate adding",
			args: args{
				userName: "us1",
				password: "pass1",
			},
			isErr:     true,
			errString: "23505",
		},
	}
	ctx = context.Background()
	dataBase, err := ConnectUsersTable(ctx, dbEndPoint)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}
	err = dataBase.UsersTableCreation(ctx, testTableName)
	if err != nil {
		fmt.Printf("error  table creation %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dataBase.AddUser(ctx, testTableName, tt.args.userName, tt.args.password)
			assert.Equal(t, tt.isErr, err != nil)
			if err != nil {
				assert.ErrorContains(t, err, tt.errString)
			}
		})
	}

	tt := tests[0]
	t.Run("correct password", func(t *testing.T) {
		err := dataBase.CheckUserPassword(ctx, testTableName, tt.args.userName, tt.args.password)
		assert.Equal(t, tt.isErr, err != nil)
		if err != nil {
			assert.ErrorContains(t, err, tt.errString)
		}
	})

	//	tt = tests[0]
	t.Run("wrong password", func(t *testing.T) {
		err := dataBase.CheckUserPassword(ctx, testTableName, tt.args.userName, tt.args.password+"a")
		assert.Equal(t, tt.isErr, err == nil)
		if err != nil {
			assert.ErrorContains(t, err, "password not match")
		}
	})
	t.Run("Right User", func(t *testing.T) {
		err := dataBase.IfUserExists(ctx, testTableName, tt.args.userName)
		assert.Equal(t, tt.isErr, err != nil)
		if err != nil {
			assert.ErrorContains(t, err, "QueryRow, error is")
		}
	})
	t.Run("Wrong User", func(t *testing.T) {
		err := dataBase.IfUserExists(ctx, testTableName, tt.args.userName+"a")
		assert.Equal(t, tt.isErr, err == nil)
		if err != nil {
			assert.ErrorContains(t, err, "QueryRow, error is")
		}
	})

	dropOrder := "DROP TABLE " + testTableName + " ;"
	tag, err := dataBase.DB.Exec(ctx, dropOrder)
	if err != nil {
		fmt.Printf("error DROP users table. Tag is \"%s\" error is %v", tag.String(), err)
		return
	}
}
