package main

// Basic imports
import (
	"log"

	"github.com/stretchr/testify/assert"
)

func (suite *TS) TestExample3() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	log.Println("testexample3")
	suite.Equal(5, suite.VariableThatShouldStartAtFive)
}
func (suite *TS) TestExample2() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	log.Println("testexample2")
	suite.Equal(5, suite.VariableThatShouldStartAtFive)
}
