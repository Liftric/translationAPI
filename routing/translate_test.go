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
}
