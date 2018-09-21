package routing

import (
	"github.com/gin-gonic/gin"
	"preventis.io/translationApi/model"
)

func setupTestEnvironment() *gin.Engine {
	db = model.StartDB("sqlite3", ":memory:")

	eng := model.Language{IsoCode: "en", Name: "English"}
	ger := model.Language{IsoCode: "de", Name: "German"}
	es := model.Language{IsoCode: "es", Name: "Spanish"}
	db.Create(&eng)
	db.Create(&ger)
	db.Create(&es)

	proj1 := model.Project{Name: "Shared", BaseLanguage: eng, Languages: []model.Language{ger, eng}}
	proj2 := model.Project{Name: "Base", BaseLanguage: ger, Languages: []model.Language{ger}}
	archivedProj := model.Project{Name: "Archived", BaseLanguage: ger, Archived: true}
	db.Create(&proj1)
	db.Create(&proj2)
	db.Create(&archivedProj)

	key1 := model.StringIdentifier{Identifier: "key1", Project: proj1}
	key2 := model.StringIdentifier{Identifier: "key2", Project: proj2}
	db.Create(&key1)
	db.Create(&key2)

	translation1 := model.Translation{Translation: "translation1", Identifier: key1, Language: ger}
	translation2 := model.Translation{Translation: "translation2", Identifier: key2, Language: ger, Approved: true}

	db.Create(&translation1)
	db.Create(&translation2)

	router := setupRouter()
	return router
}
