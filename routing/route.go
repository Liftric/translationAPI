package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"preventis.io/translationApi/model"
)

var db *gorm.DB

func StartRouter(database *gorm.DB) {
	db = database

	r := gin.Default()
	// display all projects (with statistics)
	r.GET("/projects", getAllActiveProjects)
	// display all archived projects
	r.GET("/projects/archived", getAllArchivedProjects)
	// display all strings, statistics and config of a project
	r.GET("/project/:id", getProject)
	// create project
	r.PUT("/project", createProject)
	// create language
	r.PUT("/language", createLanguage)
	// create key in project
	r.PUT("/key", createKey)
	// change key
	r.POST("/key/:id", updateKey)
	// create translation
	r.PUT("/translation", createTranslation)
	// change translation
	r.POST("/translation/update/:id", updateTranslation)
	// set revised for translation in a language
	r.POST("/translation/approve/:id", setApproved)
	// change project name
	r.POST("/projectName", renameProject)
	// move key to another project
	r.POST("/key/:id/move/:projectId", moveKey)
	// add language to project
	r.PUT("/project/:id/languages", addLanguageToProject)
	// set base language of project
	r.POST("/project/:id/baseLanguage", setBaseLanguage)
	// archive project
	r.DELETE("/project/:id", archiveProject)
	// delete key
	r.DELETE("/key/:id", deleteKey)
	// diff iOS strings file and db
	r.POST("/project/:id/ios", diffIOS)
	// diff android strings file and db
	r.POST("/project/:id/android", diffAndroid)
	// diff excel file and db
	r.POST("/project/:id/excel", diffExcel)
	// export ios strings
	r.GET("/project/:id/ios", exportIOS)
	// export android strings
	r.GET("/project/:id/android", exportAndroid)
	// export to excel
	r.GET("/project/:id/excel", exportExcel)

	r.Run() // listen and serve on 0.0.0.0:8080
}

type languageValidation struct {
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	IsoCode string `form:"languageCode" json:"languageCode" xml:"languageCode"  binding:"required"`
}

func createLanguage(c *gin.Context) {
	var language model.Language
	var json languageValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	language.Name = json.Name
	language.IsoCode = json.IsoCode

	if dbc := db.Create(&language); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.JSON(201, language)
}
