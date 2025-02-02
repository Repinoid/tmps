package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"oppa/internal/securitate"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_UserRegister(t *testing.T) {
	type logos struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}
	type want struct {
		code int
		//	response    string
		contentType string
		noMarshErr  bool
	}
	tests := []struct {
		testName string
		urla     string
		userName string
		password string

		want want
	}{
		{
			testName: "Right case",
			urla:     "/api/user/register",
			userName: "us1",
			password: "pass1",
			want: want{
				code:        http.StatusOK,
				noMarshErr:  true,
				contentType: "application/json",
			},
		},
		{
			testName: "User already exists",
			urla:     "/api/user/register",
			userName: "us1",
			password: "pass1",
			want: want{
				code:        http.StatusConflict, // 409 — логин уже занят;
				noMarshErr:  false,
				contentType: "application/json",
			},
		},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	securitate.UsersTable = "tutab"
	securitate.TokensTable = "tttab"
	securitate.OrdersTable = "totab"

	ctx = context.Background()
	//	var err error
	DB, err = securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			lo, _ := json.Marshal(logos{UserName: tt.userName, Password: tt.password})
			request := httptest.NewRequest(http.MethodPost, tt.urla, bytes.NewBuffer(lo))
			w := httptest.NewRecorder()
			registerUser(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				tok := struct {
					Token string
					Until time.Time
				}{}
				err = json.Unmarshal([]byte(resBody), &tok)
				require.NoError(t, err)
				require.NotEqual(t, tok.Token, "")

				//	assert.JSONEq(t, tt.want.response, string(resBody))
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})

	}
	// for _, tab := range []string{securitate.OrdersTable, securitate.TokensTable, securitate.UsersTable} {
	// 	dropOrder := "DROP TABLE " + tab + " ;"
	// 	tag, err := DB.DB.Exec(ctx, dropOrder)
	// 	if err != nil {
	// 		fmt.Printf("error DROP users table. Tag is \"%s\" error is %v", tag.String(), err)
	// 		return
	// 	}
	// }
}
func Test_UserLogin(t *testing.T) {
	type logos struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}
	type want struct {
		code int
		//	response    string
		contentType string
		noMarshErr  bool
	}
	tests := []struct {
		testName string
		urla     string
		userName string
		password string

		want want
	}{
		{
			testName: "Right case",
			urla:     "/api/user/register",
			userName: "us1",
			password: "pass1",
			want: want{
				code:        http.StatusOK,
				noMarshErr:  true,
				contentType: "application/json",
			},
		},
		{
			testName: "Wrong user",
			urla:     "/api/user/register",
			userName: "us11",
			password: "pass1",
			want: want{
				code:        http.StatusUnauthorized, // 401 — неверная пара логин/пароль;
				noMarshErr:  false,
				contentType: "application/json",
			},
		},
		{
			testName: "Wrong password",
			urla:     "/api/user/register",
			userName: "us1",
			password: "pass2",
			want: want{
				code:        http.StatusUnauthorized, // 401 — неверная пара логин/пароль;
				noMarshErr:  false,
				contentType: "application/json",
			},
		},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	securitate.UsersTable = "tutab"
	securitate.TokensTable = "tttab"
	securitate.OrdersTable = "totab"

	ctx = context.Background()
	//	var err error
	DB, err = securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			lo, _ := json.Marshal(logos{UserName: tt.userName, Password: tt.password})
			request := httptest.NewRequest(http.MethodPost, tt.urla, bytes.NewBuffer(lo))
			w := httptest.NewRecorder()
			loginUser(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				tok := struct {
					Token string
					Until time.Time
				}{}
				err = json.Unmarshal([]byte(resBody), &tok)
				require.NoError(t, err)
				require.NotEqual(t, tok.Token, "")

				//	assert.JSONEq(t, tt.want.response, string(resBody))
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})

	}
	for _, tab := range []string{securitate.OrdersTable, securitate.TokensTable, securitate.UsersTable} {
		dropOrder := "DROP TABLE " + tab + " ;"
		tag, err := DB.DB.Exec(ctx, dropOrder)
		if err != nil {
			fmt.Printf("error DROP users table. Tag is \"%s\" error is %v", tag.String(), err)
			return
		}
	}
}
