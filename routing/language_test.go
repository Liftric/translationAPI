package routing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguagesRoute(t *testing.T) {
	router := setupTestRouter()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/languages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCreateLanguageRoute(t *testing.T) {
	router := setupTestRouter()
	defer db.Close()

	var jsonStr = []byte(`{"languageCode": "de", "name": "German"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/language", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, `{"IsoCode":"de","Name":"German"}`, w.Body.String())
}
