package routing

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"preventis.io/translationApi/model"
)

func diffIOS(c *gin.Context) {
	// TODO
}

type diffDTO struct {
	Identifier string
	Create     bool
	Update     bool
}

func diffAndroid(c *gin.Context) {
	id := c.Param("id")
	lang := c.Param("lang")
	var resource androidResource
	if err := c.ShouldBindWith(&resource, binding.XML); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		if !containsLanguage(lang, project) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Language not present in project"})
			println("language not in project")
			return
		}
		var diffs []diffDTO
		for _, e := range resource.Strings {
			check := checkIdtentifier(e.Identifier, e.Translation, lang, project)
			create := false
			update := false
			if check == 1 {
				create = true
			} else if check == 2 {
				update = true
			}
			diffs = append(diffs, diffDTO{Identifier: e.Identifier, Create: create, Update: update})
		}
		c.JSON(http.StatusOK, diffs)
	}
}

func checkIdtentifier(identifier string, translation string, lang string, project model.Project) int {
	for _, i := range project.Identifiers {
		if i.Identifier == identifier {
			for _, t := range i.Translations {
				if t.LanguageRefer == lang {
					if t.Translation == translation {
						return 0
					} else {
						return 2
					}
				}
			}
		}
	}
	return 1
}

func diffExcel(c *gin.Context) {
	// TODO
}
func exportIOS(c *gin.Context) {
	id := c.Param("id")
	lang := c.Param("lang")
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		if !containsLanguage(lang, project) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Language not present in project"})
			println("language not in project")
			return
		}
		resList := getAndroidExportStrings(project, lang)
		res := ""
		for _, e := range resList {
			res += fmt.Sprintf(`"%s" = "%s";\n`, e.Identifier, e.Translation)
		}

		// not using c.XML because of formatting
		c.String(http.StatusOK, res)
	}
}

type androidResource struct {
	XMLName xml.Name        `xml:"resources" binding:"required"`
	Strings []androidString `xml:"string"`
}

type androidString struct {
	Identifier  string `xml:"name,attr"`
	Translation string `xml:",chardata"`
}

func exportAndroid(c *gin.Context) {
	id := c.Param("id")
	lang := c.Param("lang")
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		if !containsLanguage(lang, project) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Language not present in project"})
			println("language not in project")
			return
		}
		resList := getAndroidExportStrings(project, lang)
		resource := androidResource{Strings: resList}
		w := &bytes.Buffer{}
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		if err := enc.Encode(resource); err != nil {
			panic(err)
		}

		// not using c.XML because of formatting
		c.String(http.StatusOK, w.String())
	}
}

func getAndroidExportStrings(project model.Project, lang string) []androidString {
	var resList []androidString
	for _, e := range project.Identifiers {
		for _, t := range e.Translations {
			if t.LanguageRefer == lang {
				translation := androidString{Identifier: e.Identifier, Translation: t.Translation}
				resList = append(resList, translation)
			}
		}
	}
	return resList
}

func exportExcel(c *gin.Context) {
	// TODO
}
