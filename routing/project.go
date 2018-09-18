package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
