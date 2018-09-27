package routing

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"preventis.io/translationApi/model"
	"testing"
)

func TestProjectsRoute(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"Id":1,"Name":"Shared","BaseLanguage":{"IsoCode":"en","Name":"English"},"Languages":[{"IsoCode":"de","Name":"German"},{"IsoCode":"en","Name":"English"}]},{"Id":2,"Name":"Base","BaseLanguage":{"IsoCode":"de","Name":"German"},"Languages":[{"IsoCode":"de","Name":"German"}]}]`, w.Body.String())
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

func TestGetProject(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"Id":1,"Name":"Shared","BaseLanguage":{"IsoCode":"en","Name":"English"},"Languages":[{"IsoCode":"de","Name":"German"},{"IsoCode":"en","Name":"English"}],"Identifiers":[{"Id":1,"Identifier":"key1","Translations":[{"Translation":"translation1","Language":"de","Approved":false,"ImprovementNeeded":false}]},{"Id":2,"Identifier":"key2","Translations":[{"Translation":"translation2","Language":"de","Approved":false,"ImprovementNeeded":false}]}]}`, w.Body.String())
}

func TestCreateProject(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var projects []model.Project
	db.Where("Name = ?", "newProject").Find(&projects)
	assert.Empty(t, projects)

	var jsonStr = []byte(`{"name": "newProject", "baseLanguageCode": "en"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/project", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	db.Where("Name = ?", "newProject").Find(&projects)
	assert.Equal(t, 1, len(projects))

	var project model.Project
	db.Where("Name = ?", "newProject").Preload("Languages").First(&project)
	assert.Equal(t, "en", project.BaseLanguageRefer)
	assert.Equal(t, "en", project.Languages[0].IsoCode)
	assert.Equal(t, 1, len(project.Languages))

	var jsonStr2 = []byte(`{"name": "newProject2", "baseLanguageCode": "xy"}`)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/project", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/project", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 409, w3.Code)
	db.Where("Name = ?", "newProject").Find(&projects)
	assert.Equal(t, 1, len(projects))
}

func TestRenameProject(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var project model.Project
	db.Where("id = ?", 1).First(&project)

	assert.NotEqual(t, "updatedProject", project.Name)

	var jsonStr = []byte(`{"name": "updatedProject"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/name", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	db.Where("id = ?", 1).First(&project)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "updatedProject", project.Name)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/100/name", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 404, w2.Code)
}

func TestArchiveProject(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var project model.Project
	db.First(&project)

	assert.False(t, project.Archived)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/project/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	db.First(&project)
	assert.True(t, project.Archived)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/project/1", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)

	db.First(&project)
	assert.True(t, project.Archived)

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("DELETE", "/project/100", nil)
	router.ServeHTTP(w4, req4)

	assert.Equal(t, 404, w4.Code)
}

func TestAddLanguageToProject(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var project model.Project
	db.Preload("Languages").First(&project)
	assert.False(t, containsLanguage("es", project))

	var jsonStr = []byte(`{"languageCode": "es"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/languages", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var project2 model.Project
	db.Preload("Languages").First(&project2)
	assert.True(t, containsLanguage("es", project2))

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/1/languages", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 409, w2.Code)
}

func TestSetBaseLanguage(t *testing.T) {
	router := setupTestEnvironment()
	defer db.Close()

	var project model.Project
	db.First(&project)

	assert.Equal(t, "en", project.BaseLanguageRefer)

	var jsonStr = []byte(`{"languageCode": "de"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/1/baseLanguage", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	db.First(&project)
	assert.Equal(t, 200, w.Code)

	db.First(&project)

	assert.Equal(t, "de", project.BaseLanguageRefer)

	var jsonStr2 = []byte(`{"languageCode": "xy"}`)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/project/1/baseLanguage", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w2, req2)

	db.First(&project)
	assert.Equal(t, 404, w2.Code)
	assert.Equal(t, "de", project.BaseLanguageRefer)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/project/100/baseLanguage", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 404, w3.Code)
}
