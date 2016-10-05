package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

var T = &t{
	templates: template.Must(template.ParseGlob("templates/*.html")),
}

type t struct {
	templates *template.Template
}

func (t *t) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//return t.templates.ExecuteTemplate(w, name, data)
	return template.Must(template.ParseGlob("templates/*.html")).ExecuteTemplate(w, name, data)
}
