package main

import (
	"context"
	"crypto/rand"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var (
	Sugar    *zap.SugaredLogger
	baseMass []byte
)
//const Long = 33333

type TstHeapSort struct {
	suite.Suite
	t   time.Time
	ctx context.Context
}

func (suite *TstHeapSort) SetupSuite() { // выполняется перед тестами
	suite.ctx = context.Background()
	suite.t = time.Now()

	baseMass = make([]byte, Long)
	rand.Read(baseMass)

	log.Println("SetupTest() ---------------------")
}

func (suite *TstHeapSort) TearDownSuite() { // // выполняется после всех тестов
	log.Printf("Spent %v\n", time.Since(suite.t))
}

func TestHeapSortSuite(t *testing.T) {
	testHeap := new(TstHeapSort)
	testHeap.ctx = context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	Sugar = logger.Sugar()

	log.Println("before run ")
	suite.Run(t, testHeap)

}
