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

	var keys []model.StringIdentifier
	db.Where("Identifier = ? AND project_id = ?", "newIdentifier", 1).Find(&keys)
	assert.Empty(t, keys)

	var jsonStr = []byte(`{"projectId": 1, "identifier": "newIdentifier"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	db.Where("Identifier = ? AND project_id = ?", "newIdentifier", 1).Find(&keys)
	assert.NotEmpty(t, keys)

	var jsonStr2 = []byte(`{"projectId": 100, "identifier": "newIdentifier"}`)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/identifier", bytes.NewBuffer(jsonStr))
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
	req2, _ := http.NewRequest("POST", "/identifier/100", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)
}

func TestCreateTranslation(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var translations []model.Translation
	db.Where("string_identifier_id = ? AND language_refer = ?", 1, "en").Find(&translations)
	assert.Empty(t, translations)

	var jsonStr = []byte(`{"keyId": 1, "translation": "testTranslation", "languageCode": "en"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/translation", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, `{"Id":4,"Translation":"testTranslation","Language":"en","Approved":false,"ImprovementNeeded":false}`, w.Body.String())
	var translation model.Translation
	db.Where("string_identifier_id = ? AND language_refer = ?", 1, "en").First(&translation)
	assert.Equal(t, "testTranslation", translation.Translation)

	req2, _ := http.NewRequest("PUT", "/translation", bytes.NewBuffer(jsonStr))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, `{"Id":4,"Translation":"testTranslation","Language":"en","Approved":false,"ImprovementNeeded":false}`, w2.Body.String())

	db.Where("string_identifier_id = ? AND language_refer = ?", 1, "en").Find(&translations)
	assert.Equal(t, 1, len(translations))
}

func TestApproveTranslation(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var translation model.Translation
	db.Where("id = ?", 1).First(&translation)

	assert.False(t, translation.Approved)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/translation/approve/1", nil)
	router.ServeHTTP(w, req)

	db.Where("id = ?", 1).First(&translation)

	assert.Equal(t, 200, w.Code)
	assert.True(t, translation.Approved)

	// test if approval gets revoked after changing translation
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/translation/approve/1", nil)
	router.ServeHTTP(w2, req2)

	db.Where("id = ?", 1).First(&translation)
	assert.True(t, translation.Approved)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/translation/approve/100", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 404, w3.Code)
}

func TestToggleImprovementNeeded(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var translation model.Translation
	db.Where("id = ?", 1).First(&translation)

	assert.False(t, translation.ImprovementNeeded)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/translation/improvement/1", nil)
	router.ServeHTTP(w, req)

	db.Where("id = ?", 1).First(&translation)

	assert.Equal(t, 200, w.Code)
	assert.True(t, translation.ImprovementNeeded)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/translation/improvement/1", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)

	db.Where("id = ?", 1).First(&translation)
	assert.False(t, translation.ImprovementNeeded)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/translation/improvement/100", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 404, w3.Code)
}

func TestMoveKey(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var key model.StringIdentifier
	db.First(&key)

	assert.Equal(t, uint(1), key.ProjectID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/identifier/1/move/2", nil)
	router.ServeHTTP(w, req)

	db.First(&key)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, uint(2), key.ProjectID)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/identifier/1/move/2", nil)
	router.ServeHTTP(w2, req2)

	db.First(&key)
	assert.Equal(t, 304, w2.Code)
	assert.Equal(t, uint(2), key.ProjectID)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/identifier/100/move/2", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 404, w3.Code)

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("POST", "/identifier/1/move/200", nil)
	router.ServeHTTP(w4, req4)

	assert.Equal(t, 404, w4.Code)

	w5 := httptest.NewRecorder()
	req5, _ := http.NewRequest("POST", "/identifier/2/move/2", nil)
	router.ServeHTTP(w5, req5)

	assert.Equal(t, 409, w5.Code)
}

func TestDeleteKey(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var key model.StringIdentifier
	db.First(&key)
	assert.Equal(t, uint(1), key.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/identifier/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var key2 model.StringIdentifier
	db.First(&key2)
	assert.NotEqual(t, uint(1), key2.ID)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/identifier/1", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("DELETE", "/identifier/100", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 404, w3.Code)
}
