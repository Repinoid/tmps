package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oppa/internal/models"
	"oppa/internal/rual"
	"oppa/internal/securitate"
	"strconv"
	"strings"
	"time"

	"github.com/theplant/luhn"
)


var DataBase = securitate.DataBase
var ctx = models.Ctx


func RegisterUser(rwr http.ResponseWriter, req *http.Request) {

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debug("not application/json\n")
		return
	}

	rwr.Header().Set("Content-Type", "application/json")

	Token, err := securitate.BuildJWTString("someID", []byte(securitate.SECRET_KEY))
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("BuildJWTString %+v\n", err)
		return
	}

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("io.ReadAll %+v\n", err)
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
		models.Sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)
		return
	}
	models.Sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)

	err = DataBase.IfUserExists(ctx, logos.UserName)
	if err == nil {
		fmt.Printf("User exists %v\n", err)
		rwr.WriteHeader(http.StatusConflict) // 409 — логин уже занят;
		fmt.Fprintf(rwr, `{"status":"StatusConflict"}`)
		return
	}
	err = DataBase.AddUser(ctx, logos.UserName, logos.Password, Token)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest) // 400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debugf("addUser %+v err %+v\n", logos, err)
		return
	}
	tok := struct {
		Token string
		Until time.Time
	}{Token: Token, Until: time.Now().Add(securitate.TOKEN_EXP)}
	rwr.WriteHeader(http.StatusOK) // 200 — пользователь успешно зарегистрирован и аутентифицирован;
	json.NewEncoder(rwr).Encode(tok)
}

func LoginUser(rwr http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса;
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debug("not application/json\n")
		return
	}
	rwr.Header().Set("Content-Type", "application/json")

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("io.ReadAll %+v\n", err)
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
		models.Sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)
		return
	}
	models.Sugar.Debugf("json.Unmarshal %+v err %+v\n", logos, err)

	err = DataBase.IfUserExists(ctx, logos.UserName)
	if err != nil {
		fmt.Printf("User does NOT exist ERR %v\n", err)
		rwr.WriteHeader(http.StatusUnauthorized) // 401 — неверная пара логин/пароль;
		fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`)
		return
	}
	err = DataBase.CheckUserPassword(ctx, logos.UserName, logos.Password)
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
	err = DataBase.UpdateToken(ctx, logos.UserName, Token)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("UpdateToken %+v\n", err)
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
		models.Sugar.Debug("not text/plain \n")
		return
	}
	tokenStr := req.Header.Get("Authorization")
	tokenStr, niceP := strings.CutPrefix(tokenStr, "Bearer <") // обрезаем -- Bearer <token>
	tokenStr, niceS := strings.CutSuffix(tokenStr, ">")

	var tokenID int64
	//	err := DataBase.GetIDByToken(ctx, tokenStr, &tokenID)	// получаем ID пользователя по полученному токену

	if (!niceP) || (!niceS) || (DataBase.GetIDByToken(ctx, tokenStr, &tokenID) != nil) { // если неверная строка в Authorization - до GetIDByToken дело не дойдёт
		rwr.WriteHeader(http.StatusUnauthorized)            // 401 — неверная пара логин/пароль;
		fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`) // либо токена неверный формат, либо по нему нет юзера в базе
		models.Sugar.Debug("Authorization header\n")
		return
	}

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()

	orderStr := string(telo)                            // telo - []byte. В нём только номер заказа
	orderNum, err := strconv.ParseInt(orderStr, 10, 64) //
	if err != nil || (!luhn.Valid(int(orderNum))) {     // если не распарсилось или не по ЛУНУ
		rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
		fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
		models.Sugar.Debug("ordernum err\n")
		return
	}
	var orderID int64
	err = DataBase.GetIDByOrder(ctx, orderNum, &orderID)
	if err != nil { // если такого номера заказа нет в базе вносим его

		orderStat, statusCode := rual.GetFromAccrual(orderStr)

		//err =  // tokenID)	- ID пользователя по полученному токену
		if statusCode != http.StatusOK || DataBase.UpLoadOrderByID(ctx, tokenID, orderNum, orderStat.Status, orderStat.Accrual) != nil {
			rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
			fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
			models.Sugar.Debug("ordernum err\n")
			return
		}
		rwr.WriteHeader(http.StatusAccepted) //202 — новый номер заказа принят в обработку;
		fmt.Fprintf(rwr, `{"status":"StatusAccepted"}`)
		return
	}
	if orderID == tokenID {
		rwr.WriteHeader(http.StatusOK) // 200 — номер заказа уже был загружен ЭТИМ пользователем;
		fmt.Fprintf(rwr, `{"status":"StatusOK"}`)
		models.Sugar.Debug("ordernum err\n")
		return
	}
	rwr.WriteHeader(http.StatusConflict) // 409 — номер заказа уже был загружен другим пользователем;
	fmt.Fprintf(rwr, `{"status":"StatusConflict"}`)
	models.Sugar.Debug("ordernum err\n")
}
