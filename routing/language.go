package routing

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

func getLanguages(c *gin.Context) {
	var languages []model.Language
	db.Find(&languages)
	c.JSON(200, languages)
}

type languageValidation struct {
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	IsoCode string `form:"languageCode" json:"languageCode" xml:"languageCode"  binding:"required"`
}

func createLanguage(c *gin.Context) {
	var json languageValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var language model.Language
	language.Name = json.Name
	language.IsoCode = json.IsoCode
	if dbc := db.Create(&language); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.JSON(201, language)
}
