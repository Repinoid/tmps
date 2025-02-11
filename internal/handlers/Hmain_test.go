package handlers

import (
	"context"
	"fmt"
	"log"
	"oppa/internal/models"
	"oppa/internal/securitate"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type TstHandlers struct {
	suite.Suite
}

func (suite *TstHandlers) SetupSuite() {
	//var err error
	ctx := context.Background()
	dataBase, err := securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return
	}
	defer dataBase.DB.Close(ctx)
	for _, tab := range []string{"orders", "tokens", "withdrawn", "accounts"} {
		dropOrder := "DROP TABLE " + tab + " ;"
		_, err := dataBase.DB.Exec(ctx, dropOrder)
		if err != nil {
			log.Printf("error DROP %s table. %v", tab, err)
			return
		}
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	models.Sugar = *logger.Sugar()
	log.Println("SetupTest() ---------------------")
}

func (suite *TstHandlers) TearDownSuite() {
}

//	func (suite *TSuite) BeforeTest(suiteName, testName string) {
//		log.Println("BeforeTest()", suiteName, testName)
//	}
//
//	func (suite *TSuite) AfterTest(suiteName, testName string) {
//		log.Println("AfterTest()", suiteName, testName)
//	}
func TestHandlersSuite(t *testing.T) {
	log.Println("before run")
	suite.Run(t, new(TstHandlers))
}
