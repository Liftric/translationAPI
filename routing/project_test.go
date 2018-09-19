package routing

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectsRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Id":1,"Name":"Shared","BaseLanguage":{"IsoCode":"en","Name":"English"},"Languages":[{"IsoCode":"de","Name":"German"}]},{"Id":2,"Name":"Base","BaseLanguage":{"IsoCode":"de","Name":"German"},"Languages":[]}]`, w.Body.String())
}

func TestArchivedProjectsRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/archived", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Id":3,"Name":"Archived","BaseLanguage":{"IsoCode":"","Name":""},"Languages":[]}]`, w.Body.String())
}
