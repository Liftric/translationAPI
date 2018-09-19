package routing

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"preventis.io/translationApi/model"
	"testing"
)

func TestCreateIdentifierRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var jsonStr = []byte(`{"projectId": 1, "identifier": "newIdentifier"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var key model.StringIdentifier
	errors := db.Where("Identifier = ? AND project_id = ?", "newIdentifier", 1).First(&key).GetErrors()

	assert.Empty(t, errors)

	var jsonStr2 = []byte(`{"projectId": 100, "identifier": "newIdentifier"}`)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)

	var jsonStr3 = []byte(`{"projectId": 1, "identifier": "newIdentifier"}`)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr3))
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 409, w3.Code)
}

func TestUpdateIdentifier(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var key model.StringIdentifier
	db.Where("id = ?", 1).First(&key)

	assert.NotEqual(t, "updatedIdentifier", key.Identifier)

	var jsonStr = []byte(`{"identifier": "updatedIdentifier"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/identifier/1", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	db.Where("id = ?", 1).First(&key)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "updatedIdentifier", key.Identifier)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/identifier/100", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)
}
