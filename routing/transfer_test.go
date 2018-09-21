package routing

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExportAndroid(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1/android/de", nil)
	router.ServeHTTP(w, req)

	println(w.Body.String())

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<resources>\n  <string name=\"key1\">translation1</string>\n</resources>", w.Body.String())
}

func TestExportIos(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1/ios/de", nil)
	router.ServeHTTP(w, req)

	println(w.Body.String())

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "\"key1\" = \"translation1\";\\n", w.Body.String())
}
