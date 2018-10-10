package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

type stringIdentifierValidation struct {
	ProjectId  int    `form:"projectId" json:"projectId" xml:"projectId"  binding:"required"`
	Identifier string `form:"identifier" json:"identifier" xml:"identifier"  binding:"required"`
}

func createIdentifier(c *gin.Context) {
	var json stringIdentifierValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	key.Identifier = json.Identifier
	var project model.Project
	if err := db.Where("id = ?", json.ProjectId).
		Preload("Languages").
		Preload("BaseLanguage").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Project = project

	var existingKeys []model.StringIdentifier
	db.Where("Identifier = ? AND project_id = ?", "newIdentifier", 1).Find(&existingKeys)
	if len(existingKeys) > 0 {
		c.Status(409)
		return
	}

	if dbc := db.Create(&key); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}

	c.JSON(201, identifierToDTO(key))
}

func identifierToDTO(key model.StringIdentifier) identifierDTO {
	translations := []translationDTO{}
	for _, t := range key.Translations {
		translation := translationDTO{Id: t.ID, Translation: t.Translation, Language: t.LanguageRefer, Approved: t.Approved, ImprovementNeeded: t.ImprovementNeeded}
		translations = append(translations, translation)
	}
	return identifierDTO{Id: key.ID, Identifier: key.Identifier, Translations: translations}
}

type updateIdentifierValidation struct {
	Key string `form:"identifier" json:"identifier" xml:"identifier"  binding:"required"`
}

func updateIdentifier(c *gin.Context) {
	id := c.Param("id")
	var json updateIdentifierValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).Preload("Translations").First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Identifier = json.Key
	db.Save(&key)
	c.JSON(200, identifierToDTO(key))
}

type translationValidation struct {
	KeyId       int    `form:"keyId" json:"keyId" xml:"keyId"  binding:"required"`
	Translation string `form:"translation" json:"translation" xml:"translation"`
	Language    string `form:"languageCode" json:"languageCode" xml:"languageCode"  binding:"required"`
}

func upsertTranslation(c *gin.Context) {
	var json translationValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var translations []model.Translation
	db.Where("string_identifier_id = ? AND language_refer = ?", json.KeyId, json.Language).Preload("Language").Find(&translations)
	if len(translations) > 0 {
		translation := translations[0]
		if translation.Translation == json.Translation {
			translationDTO := translationDTO{Id: translation.ID, Translation: translation.Translation, Language: translation.Language.IsoCode, Approved: translation.Approved, ImprovementNeeded: translation.ImprovementNeeded}
			c.JSON(http.StatusOK, translationDTO)
			return
		}
		translation.Translation = json.Translation
		translation.Approved = false
		translation.ImprovementNeeded = false
		db.Save(&translation)
		revision := model.Revision{RevisionTranslation: translation.Translation, Approved: translation.Approved, Translation: translation}
		db.Create(&revision)
		translationDTO := translationDTO{Id: translation.ID, Translation: translation.Translation, Language: translation.Language.IsoCode, Approved: translation.Approved, ImprovementNeeded: translation.ImprovementNeeded}
		c.JSON(http.StatusOK, translationDTO)
		return
	}

	var translation model.Translation
	translation.Translation = json.Translation
	var key model.StringIdentifier
	if err := db.Where("id = ?", json.KeyId).Preload("Project").Preload("Project.Languages").First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	translation.Identifier = key

	var containsLang = false
	for _, e := range key.Project.Languages {
		if e.IsoCode == json.Language {
			containsLang = true
			translation.Language = e
		}
	}

	if !containsLang {
		c.AbortWithStatusJSON(404, gin.H{"error": "Project does not contain language, please add it first."})
		fmt.Println("project does not contain language")
		return
	}

	if dbc := db.Create(&translation); dbc.Error != nil {
		revision := model.Revision{RevisionTranslation: translation.Translation, Approved: translation.Approved, Translation: translation}
		db.Create(&revision)
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	translationDTO := translationDTO{Id: translation.ID, Translation: translation.Translation, Language: translation.Language.IsoCode, Approved: translation.Approved, ImprovementNeeded: translation.ImprovementNeeded}
	c.JSON(http.StatusCreated, translationDTO)
}

func setApproved(c *gin.Context) {
	id := c.Param("id")
	var translation model.Translation
	if err := db.Where("id = ?", id).First(&translation).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		translation.Approved = true
		db.Save(&translation)
		var revision model.Revision
		db.Where("TranslationID = ?", translation.ID).Order("CreatedAt").Last(&revision)
		revision.Approved = translation.Approved
		db.Save(&revision)
		c.Status(200)
	}
}

func toggleImprovementNeeded(c *gin.Context) {
	id := c.Param("id")
	var translation model.Translation
	if err := db.Where("id = ?", id).First(&translation).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		translation.ImprovementNeeded = !translation.ImprovementNeeded
		db.Save(&translation)
		c.Status(200)
	}
}

func moveKey(c *gin.Context) {
	id := c.Param("id")
	projectId := c.Param("projectId")

	var project model.Project
	if err := db.Where("id = ?", projectId).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Project = project
	db.Save(&key)
	c.Status(200)
}

func deleteKey(c *gin.Context) {
	id := c.Param("id")
	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	db.Delete(&key)
	c.Status(200)
}
