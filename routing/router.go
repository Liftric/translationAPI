package routing

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
)

var db *gorm.DB

func setupRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{getEnv("FRONTEND_URL", "http://localhost:3000")}
	config.AllowMethods = []string{"GET", "PUT", "POST", "DELETE"}
	r.Use(cors.New(config))

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
		projectRoutes.POST("/:id/ios/:lang", diffIOS)
		// diff android strings file and db
		projectRoutes.POST("/:id/android/:lang", diffAndroid)
		// diff excel file and db
		projectRoutes.POST("/:id/excel/:lang", diffExcel)
		// export ios strings
		projectRoutes.GET("/:id/ios/:lang", exportIOS)
		// export android strings
		projectRoutes.GET("/:id/android/:lang", exportAndroid)
		// export to excel
		projectRoutes.GET("/:id/csv", exportCsv)
		// export strings in json format
		projectRoutes.GET("/:id/json", exportJSON)
	}

	translationRoutes := r.Group("/translation")
	{
		// create or update translation
		translationRoutes.PUT("", upsertTranslation)
		// set revised for translation in a language
		translationRoutes.POST("/approve/:id", setApproved)
		// toggle improvement needed
		translationRoutes.POST("/improvement/:id", toggleImprovementNeeded)
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func StartRouter(database *gorm.DB) {
	db = database

	r := setupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
