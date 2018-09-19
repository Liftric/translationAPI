package routing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguagesRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/languages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"IsoCode":"en","Name":"English"},{"IsoCode":"de","Name":"German"}]`, w.Body.String())
}

func TestCreateLanguageRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var jsonStr = []byte(`{"languageCode": "fr", "name": "French"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/language", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, `{"IsoCode":"fr","Name":"French"}`, w.Body.String())
}
