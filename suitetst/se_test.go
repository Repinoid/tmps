package main

// Basic imports
import (
	"log"

	"github.com/stretchr/testify/assert"
)

func (suite *TS) TestExample5() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	log.Println("testexample5")
	suite.Equal(5, suite.VariableThatShouldStartAtFive)
}

func (suite *TS) TestExample1() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	log.Println("testexample1 ", suite.T().Name())
	suite.Equal(5, suite.VariableThatShouldStartAtFive)
	suite.Run("named", func() {
		log.Println("inside 1")
	})
}
