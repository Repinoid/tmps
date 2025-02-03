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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_DropTables(t *testing.T) {
	ctx = context.Background()
	//	var err error
	DB, err := securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}
	for _, tab := range []string{"orders", "tokens", "accounts"} {
		dropOrder := "DROP TABLE " + tab + " ;"
		tag, err := DB.DB.Exec(ctx, dropOrder)
		if err != nil {
			fmt.Printf("error DROP users table. Tag is \"%s\" error is %v", tag.String(), err)
			return
		}
	}
}

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
			testName: "Right case",
			urla:     "/api/user/register",
			userName: "us111",
			password: "pass1",
			want: want{
				code:        http.StatusOK,
				noMarshErr:  true,
				contentType: "application/json",
			},
		},
		{
			testName: "Right case",
			urla:     "/api/user/register",
			userName: "us222",
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
			request.Header.Set("Content-Type", "application/json")
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

				var tokenFromBase string
				err = DB.GetToken(ctx, tt.userName, &tokenFromBase)
				if err != nil {
					fmt.Printf("tst %v", err)
					return
				}
				assert.Equal(t, tok.Token, tokenFromBase, "токен из базы не равен токену из ответа хандлера")
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})

	}

}
func Test_UserLogin(t *testing.T) {
	type logos struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}
	type want struct {
		code        int
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
			request.Header.Set("Content-Type", "application/json")
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
}
func Test_PutOrder(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		testName    string
		urla        string
		userName    string
		orderNum    int
		ContentType string
		TokenSuffix string

		want want
	}{
		{
			testName:    "Right PUT",
			urla:        "/api/user/orders",
			userName:    "us1",
			orderNum:    111,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusAccepted,
				response:    `{"status":"StatusAccepted"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},
		{
			testName:    "Right PUT 222",
			urla:        "/api/user/orders",
			userName:    "us222",
			orderNum:    222,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusAccepted,
				response:    `{"status":"StatusAccepted"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},
		{
			testName:    "Already PUT",
			urla:        "/api/user/orders",
			userName:    "us1",
			orderNum:    111,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusOK,
				response:    `{"status":"StatusOK"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},
		{
			testName:    "Other PUT",
			urla:        "/api/user/orders",
			userName:    "us1",
			orderNum:    222,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusConflict,
				response:    `{"status":"StatusConflict"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},

		{
			testName:    "Wrong Content Type",
			urla:        "/api/user/orders",
			userName:    "us1",
			orderNum:    111,
			ContentType: "application/json",
			want: want{
				code:        http.StatusBadRequest,
				response:    `{"status":"StatusBadRequest"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},
		{
			testName:    "Wrong PUT user not exist",
			urla:        "/api/user/orders",
			userName:    "us10",
			orderNum:    111,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusUnauthorized,
				response:    `{"status":"StatusUnauthorized"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">",
		},
		{
			testName:    "Wrong TOKEN string",
			urla:        "/api/user/orders",
			userName:    "us1",
			orderNum:    111,
			ContentType: "text/plain",
			want: want{
				code:        http.StatusUnauthorized,
				response:    `{"status":"StatusUnauthorized"}`,
				contentType: "application/json",
			},
			TokenSuffix: ">>",
		},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	ctx = context.Background()
	//	var err error
	DB, err = securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var token string
			err = DB.GetToken(ctx, tt.userName, &token)
			tokenStr := "Bearer <" + token + tt.TokenSuffix

			request := httptest.NewRequest(http.MethodPost, tt.urla, bytes.NewBufferString(strconv.Itoa(tt.orderNum)))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", tt.ContentType)
			request.Header.Set("Authorization", tokenStr)
			PutOrder(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.JSONEq(t, tt.want.response, string(resBody))

		})

	}

}
