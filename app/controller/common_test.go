package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func Test_resposneWithData(t *testing.T) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)

	resposneWithData(c, gin.H{"testKey1": "testVal1", "testKey2": []string{"testVal2_item1", "testVal2_item2"}})

	assert.Equal(t, r.Code, http.StatusOK)
	expectedResp, _ := json.Marshal(gin.H{"success": true, "testKey1": "testVal1", "testKey2": []string{"testVal2_item1", "testVal2_item2"}})
	assert.Equal(t, string(expectedResp), r.Body.String())
}

func Test_resposneWithError(t *testing.T) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)

	respondeWithError(c, http.StatusBadRequest, errors.New("balance insufficient"))

	assert.Equal(t, r.Code, http.StatusBadRequest)
	expectedResp, _ := json.Marshal(gin.H{"success": false, "error": "balance insufficient"})
	assert.Equal(t, string(expectedResp), r.Body.String())
}
