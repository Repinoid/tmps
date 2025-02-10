package main

// Basic imports
import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TS struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TS) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}
func (suite *TS) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}
func (suite *TS) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}
func TestExampleTestSuite(t *testing.T) {
	log.Println("before run")
	suite.Run(t, new(TS))
}
