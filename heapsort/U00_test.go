package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

type TstHeapSort struct {
	suite.Suite
	t   time.Time
	ctx context.Context
}

func (suite *TstHeapSort) SetupSuite() { // выполняется перед тестами
	suite.ctx = context.Background()
	suite.t = time.Now()

	log.Println("SetupTest() ---------------------")
}

func (suite *TstHeapSort) TearDownSuite() { // // выполняется после всех тестов
	log.Printf("Spent %v\n", time.Since(suite.t))
}

func TestHandlersSuite(t *testing.T) {
	testBase := new(TstHeapSort)
	testBase.ctx = context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	Sugar = logger.Sugar()

	log.Println("before run ")
	suite.Run(t, testBase)

}
