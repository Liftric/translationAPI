package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func setupRouter() *gin.Engine {
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
		projectRoutes.POST("/:id/languages", addLanguageToProject)
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

	identifierRoutes := r.Group("/identifier")
	{
		// create key in project
		identifierRoutes.PUT("", createIdentifier)
		// change key
		identifierRoutes.POST("/:id", updateIdentifier)
		// move key to another project
		identifierRoutes.POST("/:id/move/:projectId", moveKey)
		// delete key
		identifierRoutes.DELETE("/:id", deleteKey)
	}

	// create language
	r.PUT("/language", createLanguage)
	// get languages
	r.GET("/languages", getLanguages)
	return r
}

func StartRouter(database *gorm.DB) {
	db = database

	r := setupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
