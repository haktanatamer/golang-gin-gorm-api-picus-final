package api_test

import (
	"api-gin/package/internal/api"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createMockHandlers() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api.RegisterHandlersTest(r)
	return r
}

//post method test
func TestUserAdd(t *testing.T) {
	r := createMockHandlers()
	req, err := http.NewRequest("POST", "/user/add", nil)

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

//get method test
func TestOrderList(t *testing.T) {
	r := createMockHandlers()
	req, err := http.NewRequest("GET", "/order/list", nil)

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
