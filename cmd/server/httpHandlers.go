package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oppa/internal/securitate"
	"strconv"
	"strings"
	"time"
)

func registerUser(rwr http.ResponseWriter, req *http.Request) {

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debug("not application/json\n")
		return
	}

	rwr.Header().Set("Content-Type", "application/json")

	Token, err := securitate.BuildJWTString("someID", []byte(securitate.SECRET_KEY))
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		sugar.Debugf("BuildJWTString %+v\n", err)
		return
	}

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()

	logos := struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}{}
	err = json.Unmarshal([]byte(telo), &logos)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)
		return
	}
	sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)

	err = DB.IfUserExists(ctx, logos.UserName)
	if err == nil {
		fmt.Printf("User exists %v\n", err)
		rwr.WriteHeader(http.StatusConflict) // 409 — логин уже занят;
		fmt.Fprintf(rwr, `{"status":"StatusConflict"}`)
		return
	}
	err = DB.AddUser(ctx, logos.UserName, logos.Password, Token)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debugf("addUser %+v err %+v\n", logos, err)
		return
	}
	tok := struct {
		Token string
		Until time.Time
	}{Token: Token, Until: time.Now().Add(securitate.TOKEN_EXP)}
	rwr.WriteHeader(http.StatusOK) // 200 — пользователь успешно зарегистрирован и аутентифицирован;
	json.NewEncoder(rwr).Encode(tok)
}

func loginUser(rwr http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debug("not application/json\n")
		return
	}
	rwr.Header().Set("Content-Type", "application/json")

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()

	logos := struct {
		UserName string `json:"login"`
		Password string `json:"password"`
	}{}
	err = json.Unmarshal([]byte(telo), &logos)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)
		return
	}
	sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)

	err = DB.IfUserExists(ctx, logos.UserName)
	if err != nil {
		fmt.Printf("User does NOT exist ERR %v\n", err)
		rwr.WriteHeader(http.StatusUnauthorized) // 401 — неверная пара логин/пароль;
		fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`)
		return
	}
	err = DB.CheckUserPassword(ctx, logos.UserName, logos.Password)
	if err != nil {
		fmt.Printf("Wrong password ERR %v\n", err)
		rwr.WriteHeader(http.StatusUnauthorized) // 401 — неверная пара логин/пароль;
		fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`)
		return
	}
	Token, err := securitate.BuildJWTString("someID", []byte(securitate.SECRET_KEY))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	tok := struct {
		Token string
		Until time.Time
	}{Token: Token, Until: time.Now().Add(securitate.TOKEN_EXP)}
	rwr.WriteHeader(http.StatusOK) // 200 — пользователь успешно зарегистрирован и аутентифицирован;
	json.NewEncoder(rwr).Encode(tok)
}

// --------------------------------------------------------------------------------------------------
func PutOrder(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса; не text/plain
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		sugar.Debug("not text/plain \n")
		return
	}
	tokenStr := req.Header.Get("Authorization")
	tokenStr, niceP := strings.CutPrefix(tokenStr, "Bearer <") // обрезаем -- Bearer <token>
	tokenStr, niceS := strings.CutSuffix(tokenStr, ">")

	var tokenID int64
	err := DB.GetIDByToken(ctx, tokenStr, &tokenID)

	if (!niceP) || (!niceS) || (err != nil) {
		rwr.WriteHeader(http.StatusUnauthorized) // 401 — неверная пара логин/пароль;
		fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`)
		sugar.Debug("Authorization header\n")
		return
	}

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()
	orderStr := string(telo)
	orderNum, err := strconv.ParseInt(orderStr, 10, 64)
	if err != nil {
		rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
		fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
		sugar.Debug("ordernum err\n")
		return
	}
	err = DB.UpLoadOrderByID(ctx, tokenID, orderNum)
	if err != nil {
		rwr.WriteHeader(http.StatusConflict) // 409 — номер заказа уже был загружен другим пользователем;
		fmt.Fprintf(rwr, `{"status":"StatusConflict"}`)
		sugar.Debug("ordernum err\n")
		return
	}
	rwr.WriteHeader(http.StatusAccepted) //202 — новый номер заказа принят в обработку;
	fmt.Fprintf(rwr, `{"status":"StatusAccepted"}`)

}
