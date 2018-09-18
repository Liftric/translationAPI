package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

func getAllProjects(c *gin.Context) {
	var projects []model.Project
	db.Find(&projects)
	c.JSON(200, projects)
}
func getProject(c *gin.Context) {
	id := c.Param("id")
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, project)
	}
}

type ProjectValidation struct {
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

func createProject(c *gin.Context) {
	var project model.Project
	var json ProjectValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project.Name = json.Name

	if dbc := db.Create(&project); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.JSON(200, project)
}

func renameProject(c *gin.Context) {
	// TODO
}

func archiveProject(c *gin.Context) {
	// TODO
}
