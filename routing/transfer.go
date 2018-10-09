package routing

import (
	"bytes"
	"encoding/csv"
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
	Identifier     string
	IdentifierId   uint
	Create         bool
	Update         bool
	ToChange       bool
	TranslationNew string
	TranslationOld string
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
			check, translationOld, identifierId := checkIdtentifier(e.Identifier, e.Translation, lang, project)
			create := false
			update := false
			if check == 1 {
				create = true
			} else if check == 2 {
				update = true
			}
			toChange := check > 0
			diffs = append(
				diffs,
				diffDTO{
					Identifier:     e.Identifier,
					IdentifierId:   identifierId,
					Create:         create,
					Update:         update,
					ToChange:       toChange,
					TranslationOld: translationOld,
					TranslationNew: e.Translation})
		}
		c.JSON(http.StatusOK, diffs)
	}
}

func checkIdtentifier(identifier string, translation string, lang string, project model.Project) (int, string, uint) {
	for _, i := range project.Identifiers {
		if i.Identifier == identifier {
			for _, t := range i.Translations {
				if t.LanguageRefer == lang {
					if t.Translation == translation {
						return 0, t.Translation, i.ID
					} else {
						return 2, t.Translation, i.ID
					}
				}
			}
			return 2, "", i.ID
		}
	}
	return 1, "", 0
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
			res += fmt.Sprintf(`"%s" = "%s";`, e.Identifier, e.Translation)
			res += fmt.Sprintf("\n")
		}

		c.Header("Content-Disposition", `attachment; filename="Localizable.strings"`)

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
		c.AbortWithStatus(http.StatusNotFound)
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

		c.Header("Content-Type", "text/xml")
		c.Header("Content-Disposition", `attachment; filename="strings.xml"`)

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

func exportCsv(c *gin.Context) {
	id := c.Param("id")
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
		return
	} else {
		var resList [][]string
		for _, e := range project.Identifiers {
			var line []string
			line = append(line, e.Identifier)
			for _, l := range project.Languages {
				appended := false
				for _, t := range e.Translations {
					if t.LanguageRefer == l.IsoCode {
						line = append(line, t.Translation)
						appended = true
					}
				}
				if !appended {
					line = append(line, "")
				}
			}
			resList = append(resList, line)
		}
		var header []string
		header = append(header, "Identifier")
		for _, l := range project.Languages {
			header = append(header, l.Name)
		}
		w := &bytes.Buffer{}
		enc := csv.NewWriter(w)
		enc.UseCRLF = true
		enc.Write(header)
		if err := enc.WriteAll(resList); err != nil {
			panic(err)
		}
		enc.Flush()

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", `attachment; filename="strings.csv"`)

		c.String(http.StatusOK, w.String())
	}
}
