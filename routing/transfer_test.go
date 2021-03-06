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
	assert.Equal(t, "<resources>\n  <string name=\"key1\">translation1\\&#39;</string>\n  <string name=\"key2\">&#34;translation2&#34;</string>\n</resources>", w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/project/3/android/de", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "<resources>\n  <string name=\"LineBreakKey\">This is a string with a line\\n\\nbreak</string>\n</resources>", w2.Body.String())
}

func TestImportAndroid(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var xmlStr = []byte("<resources>\n  <string name=\"key1\">translation1'</string>\n  <string name=\"key2\">newTranslation2</string>\n  <string name=\"key3\">translation3</string>\n</resources>")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/android/de", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Identifier":"key1","IdentifierId":1,"Create":false,"Update":false,"ToChange":false,"TranslationNew":"translation1'","TranslationOld":"translation1'"},{"Identifier":"key2","IdentifierId":2,"Create":false,"Update":true,"ToChange":true,"TranslationNew":"newTranslation2","TranslationOld":"\"translation2\""},{"Identifier":"key3","IdentifierId":0,"Create":true,"Update":false,"ToChange":true,"TranslationNew":"translation3","TranslationOld":""}]`, w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/1/android/es", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/project/100/android/de", bytes.NewBuffer(xmlStr))
	router.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestImportCSV(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var csvString = []byte(`key1;translation1';
key2;newTranslation2;
key3;translation3;
;;`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/excel/de", bytes.NewBuffer(csvString))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Identifier":"key1","IdentifierId":1,"Create":false,"Update":false,"ToChange":false,"TranslationNew":"translation1'","TranslationOld":"translation1'"},{"Identifier":"key2","IdentifierId":2,"Create":false,"Update":true,"ToChange":true,"TranslationNew":"newTranslation2","TranslationOld":"\"translation2\""},{"Identifier":"key3","IdentifierId":0,"Create":true,"Update":false,"ToChange":true,"TranslationNew":"translation3","TranslationOld":""}]`, w.Body.String())

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/1/excel/es", bytes.NewBuffer(csvString))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/project/100/excel/de", bytes.NewBuffer(csvString))
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
	assert.Equal(t, "\"key1\" = \"translation1\\'\";\n\"key2\" = \"\\\"translation2\\\"\";\n", w.Body.String())
}

func TestExportCsv(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1/csv", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Identifier,German,English\r\nkey1,translation1',\r\nkey2,\"\"\"translation2\"\"\",\r\n", w.Body.String())
}

func TestExportJson(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1/json", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"de\":{\"key1\":\"translation1'\",\"key2\":\"\\\"translation2\\\"\"},\"en\":{}}", w.Body.String())
}
