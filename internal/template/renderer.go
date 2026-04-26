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
	dir       string
	funcs     template.FuncMap
	layouts   *template.Template
	templates map[string]string
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

	return r, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseFiles(
		filepath.Join(r.dir, "layouts", "base.html"),
		filepath.Join(r.dir, name+".html"))

	if err != nil {
		r.handleError(w, err)
	}

	err = t.Funcs(r.funcs).Execute(w, data)

	if err != nil {
		r.handleError(w, err)
	}
}

func (r *Renderer) handleError(w http.ResponseWriter, err error) {
	fmt.Fprintf(os.Stderr, "template error: %v\n", err)
	http.Error(w, "Internal Server Error", 500)
}
