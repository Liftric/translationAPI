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

	projectsRoutes := r.Group("/projects")
	{
		// display all projects (with statistics)
		projectsRoutes.GET("", getAllActiveProjects)
		// display all archived projects
		projectsRoutes.GET("/archived", getAllArchivedProjects)
	}

	projectRoutes := r.Group("/project")
	{
		// display all strings, statistics and config of a project
		projectRoutes.GET("/:id", getProject)
		// create project
		projectRoutes.PUT("", createProject)
		// add language to project
		projectRoutes.PUT("/:id/languages", addLanguageToProject)
		// set base language of project
		projectRoutes.POST("/:id/baseLanguage", setBaseLanguage)
		// archive project
		projectRoutes.DELETE("/:id", archiveProject)
		// change project name
		projectRoutes.POST("/:id/name", renameProject)
		// diff iOS strings file and db
		projectRoutes.POST("/:id/ios", diffIOS)
		// diff android strings file and db
		projectRoutes.POST("/:id/android", diffAndroid)
		// diff excel file and db
		projectRoutes.POST("/:id/excel", diffExcel)
		// export ios strings
		projectRoutes.GET("/:id/ios", exportIOS)
		// export android strings
		projectRoutes.GET("/:id/android", exportAndroid)
		// export to excel
		projectRoutes.GET("/:id/excel", exportExcel)
	}

	translationRoutes := r.Group("/translation")
	{
		// create translation
		translationRoutes.PUT("", createTranslation)
		// change translation
		translationRoutes.POST("/update/:id", updateTranslation)
		// set revised for translation in a language
		translationRoutes.POST("/approve/:id", setApproved)
	}

	keyRoutes := r.Group("/key")
	{
		// create key in project
		keyRoutes.PUT("", createKey)
		// change key
		keyRoutes.POST("/:id", updateKey)
		// move key to another project
		keyRoutes.POST("/:id/move/:projectId", moveKey)
		// delete key
		keyRoutes.DELETE("/:id", deleteKey)
	}

	// create language
	r.PUT("/language", createLanguage)

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
