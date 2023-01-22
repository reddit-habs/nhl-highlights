package generate

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.html
var embeddedTemplates embed.FS

var templates *template.Template = template.Must(template.ParseFS(embeddedTemplates, "templates/*.html"))

type Clip struct {
	ID          int64
	GameID      int64
	MediaURL    string
	EventID     *int64
	Title       string
	Blurb       string
	Description string
}

func Clips(clips []*Clip) ([]byte, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "clips.html", map[string]any{
		"Clips": clips,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
