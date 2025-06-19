package ui

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed templates/*.html
var templateFS embed.FS
var templates *template.Template

func init() {
	var err error
	templates = template.Must(template.ParseFS(templateFS,
		"templates/main.html",
		"templates/header.html",
		"templates/body.html",
	))
	if err != nil {
		panic("failed to parse templates: " + err.Error())
	}
}

func Render(data any) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "main", data)
	return buf.String(), err
}

func RenderToFile(outputPath string, data any) error {
	content, err := Render(data)
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}
