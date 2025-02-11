package handlers

// Basic imports
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"oppa/internal/securitate"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (suite *TstHandlers) Test01UserRegister() {
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
			testName: "Right case1",
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
			testName: "Right case111",
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
			testName: "Right case222",
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

	ctx = context.Background()
	var err error
	securitate.DataBase, err = securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}
	defer securitate.DataBase.DB.Close(ctx)

	for _, tt := range tests {
		suite.Run(tt.testName, func() {
			lo, _ := json.Marshal(logos{UserName: tt.userName, Password: tt.password})
			request := httptest.NewRequest(http.MethodPost, tt.urla, bytes.NewBuffer(lo))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			RegisterUser(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.want.code, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				tok := struct {
					Token string
					Until time.Time
				}{}
				err = json.Unmarshal([]byte(resBody), &tok)
				require.NoError(suite.T(), err)
				require.NotEqual(suite.T(), tok.Token, "")

				var tokenFromBase string
				err = securitate.DataBase.GetToken(ctx, tt.userName, &tokenFromBase)
				if err != nil {
					fmt.Printf("tst %v", err)
					return
				}
				assert.Equal(suite.T(), tok.Token, tokenFromBase, "токен из базы не равен токену из ответа хандлера")
				assert.Equal(suite.T(), tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
