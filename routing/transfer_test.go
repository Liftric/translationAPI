package routing

import (
	"bytes"
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

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<resources>\n  <string name=\"key1\">translation1</string>\n  <string name=\"key2\">translation2</string>\n</resources>", w.Body.String())
}

func TestImportAndroid(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var xmlStr = []byte("<resources>\n  <string name=\"key1\">translation1</string>\n  <string name=\"key2\">newTranslation2</string>\n  <string name=\"key3\">translation3</string>\n</resources>")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/android/de", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Identifier":"key1","Create":false,"Update":false},{"Identifier":"key2","Create":false,"Update":true},{"Identifier":"key3","Create":true,"Update":false}]`, w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/1/android/es", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/project/100/android/de", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestExportIos(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1/ios/de", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "\"key1\" = \"translation1\";\\n\"key2\" = \"translation2\";\\n", w.Body.String())
}
