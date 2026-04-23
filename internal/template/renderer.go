package template

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Renderer struct {
	dir      string
	funcs    template.FuncMap
	templates *template.Template
}

func New(dir string) (*Renderer, error) {
	r := &Renderer{
		dir: dir,
		funcs: template.FuncMap{
			"formatDate": func(t time.Time) string {
				return t.Format("2006-01-02")
			},
			"formatDateTime": func(t time.Time) string {
				return t.Format("2006-01-02T15:04")
			},
			"formatDateTimeLocal": func(t time.Time) string {
				return t.Format("2006-01-02 15:04")
			},
			"formatDuration": func(d time.Duration) string {
				hours := d.Hours()
				return fmt.Sprintf("%.1f", hours)
			},
		},
	}

	if err := r.loadTemplates(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Renderer) loadTemplates() error {
	pattern := filepath.Join(r.dir, "*.html")
	templates, err := template.New("").Funcs(r.funcs).ParseGlob(pattern)
	if err != nil {
		return err
	}
	r.templates = templates
	return nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	if err := r.templates.ExecuteTemplate(w, name+".html", data); err != nil {
		fmt.Fprintf(os.Stderr, "template error: %v\n", err)
		http.Error(w, "Internal Server Error", 500)
	}
}