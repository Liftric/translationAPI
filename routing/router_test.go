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
	db.Create(&proj1)
	db.Create(&proj2)

	router := setupRouter()
	return router
}
