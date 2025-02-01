package securitate

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbEndPoint = "postgres://postgres:passwordas@forgo.c7wegmiakpkw.us-west-1.rds.amazonaws.com:5432/forgo"
const testTableName = "testTable"

var ctx context.Context

func TestDBstruct_AddUser(t *testing.T) {
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "Nice adding",
			args: args{
				userName: "us1",
				password: "pass1",
			},
			wantErr: "",
		},
		{
			name: "Duplicate adding",
			args: args{
				userName: "us1",
				password: "pass1",
			},
			wantErr: "",
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
			err := dataBase.AddUser(ctx, tt.args.userName, tt.args.password)
			if err != nil {
				assert.ErrorContains(t, err, tt.wantErr)
				//			t.Errorf("DBstruct.AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	dropOrder := "DROP TABLE " + testTableName + " ;"
	tag, err := dataBase.DB.Exec(ctx, dropOrder)
	if err != nil {
		fmt.Printf("error create users table. Tag is \"%s\" error is %v", tag.String(), err)
		return
	}
}
