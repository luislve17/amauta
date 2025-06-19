package ui

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFiles(
		filepath.Join("templates", "main.html"),
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "body.html"),
	)
	if err != nil {
		panic("failed to parse templates: " + err.Error())
	}
}

// Render renders the template to a string (testable)
func Render(data any) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "main", data)
	return buf.String(), err
}

// RenderToFile renders and writes to file
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
