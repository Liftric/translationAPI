package routing

import (
	"github.com/gin-gonic/gin"
	"preventis.io/translationApi/model"
)

func setupTestEnvironment() *gin.Engine {
	db = model.StartDB("sqlite3", ":memory:")

	eng := model.Language{IsoCode: "en", Name: "English"}
	ger := model.Language{IsoCode: "de", Name: "German"}
	db.Create(&eng)
	db.Create(&ger)

	proj1 := model.Project{Name: "Shared", BaseLanguage: eng, Languages: []model.Language{ger}}
	proj2 := model.Project{Name: "Base", BaseLanguage: ger}
	archivedProj := model.Project{Name: "Archived", BaseLanguage: ger, Archived: true}
	db.Create(&proj1)
	db.Create(&proj2)
	db.Create(&archivedProj)

	key1 := model.StringIdentifier{Identifier: "key1", Project: proj1}
	key2 := model.StringIdentifier{Identifier: "key2", Project: proj2}
	db.Create(&key1)
	db.Create(&key2)

	router := setupRouter()
	return router
}
