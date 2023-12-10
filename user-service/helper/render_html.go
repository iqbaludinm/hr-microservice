package helper

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

func RenderHTML(templateFS embed.FS, templateName string, data interface{}) (bytes.Buffer, error) {
	// Define the functions
	funcMap := template.FuncMap{
		"add": func(x int) int { return x + 1 },
	}
	// render html from template
	parsedTemplate, err := template.New(templateName).Funcs(funcMap).ParseFS(templateFS, fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		return bytes.Buffer{}, err
	}

	var html bytes.Buffer
	err = parsedTemplate.Execute(&html, data)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return html, nil
}
