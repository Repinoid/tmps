package rual

// Basic imports
import (
	"log"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func (suite *TSuite) Test01Setup() {

	for idx := range marks {
		err := LoadGood(idx, idx%5, 1000)
		assert.NoErrorf(suite.T(), err, "err %w", err)
	}
	log.Println("testexample5")

}
func (suite *TSuite) Test02GetFromAccrual() {

	for idx := range marks {
		Order := strconv.Itoa(Luhner(idx))
		orderStat, StatusCode, err := GetFromAccrual(Order)
		assert.NoErrorf(suite.T(), err, "err %v", err)
		assert.Equal(suite.T(), http.StatusOK, StatusCode)
		log.Println(orderStat)

	}
	log.Println("TestGetFromAccrual")

}
