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
func TestExampleTestSuite(t *testing.T) {
	log.Println("before run")
	suite.Run(t, new(TS))
}
