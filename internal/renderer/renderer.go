package renderer

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"maps"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/utility"
)

type Renderer struct {
	funcs              template.FuncMap
	templates          map[string]*template.Template
	defaultData        map[string]any
	staticHashCache    map[string]string
	webFS              embed.FS
	staticFilesRoot    string
	webStaticFilesRoot string
	templateSuffix     string
}

func New(webFS embed.FS, staticFilesRoot string, webStaticFilesRoot string, templateFilesRoot string, templateSuffix string, defaultData map[string]any) (*Renderer, error) {
	r := &Renderer{
		defaultData:        defaultData,
		templates:          make(map[string]*template.Template),
		staticHashCache:    make(map[string]string),
		webFS:              webFS,
		staticFilesRoot:    strings.TrimSuffix(staticFilesRoot, "/"),
		webStaticFilesRoot: strings.TrimSuffix(webStaticFilesRoot, "/"),
		templateSuffix:     templateSuffix,
	}

	r.funcs = template.FuncMap{
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
			return utility.GetFormattedDuration(d, false)
		},
		"formatDurationCompressed": func(d time.Duration) string {
			return utility.GetFormattedDuration(d, true)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"daysLater": func(start, end time.Time) int {
			zeroedStart := utility.ZeroTimeComponents(start)
			diff := end.Sub(zeroedStart)
			return int(math.Floor(diff.Hours() / 24.0))
		},
		"seq": func(start, end int) []int {
			n := end - start + 1
			if n <= 0 {
				return nil
			}
			s := make([]int, n)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"join": strings.Join,
		"isAfter": func(checkAfter, base time.Time) bool {
			return checkAfter.After(base)
		},
	}

	templateFS, err := fs.Sub(webFS, templateFilesRoot)
	if err != nil {
		return nil, err
	}

	err = r.initTemplates(templateFS)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Renderer) initTemplates(templateFS fs.FS) error {
	templateSuffix := r.templateSuffix
	baseTemplateName := filepath.Join("base" + templateSuffix)
	partials, err := fs.Glob(templateFS, filepath.Join("partials", "*"+templateSuffix))
	if err != nil {
		return err
	}
	baseFiles := append([]string{baseTemplateName}, partials...)

	foundLayoutFiles, err := fs.Glob(templateFS, filepath.Join("layouts", "*"+templateSuffix))
	if err != nil {
		return err
	}

	slog.Debug("Parsing base template", "baseTemplateName", baseTemplateName, "baseFiles", baseFiles)
	baseTemplate, err := template.New(baseTemplateName).Funcs(r.funcs).ParseFS(templateFS, baseFiles...)
	if err != nil {
		return err
	}
	// https://stackoverflow.com/questions/50842389/parsing-multiple-templates-in-go
	for _, layoutFile := range foundLayoutFiles {
		layoutName := filepath.Base(layoutFile)
		slog.Debug("Parsing layout", "templateName", layoutName)

		templ, err := baseTemplate.Clone()
		if err != nil {
			return err
		}
		templ, err = templ.ParseFS(templateFS, layoutFile)
		if err != nil {
			return err
		}

		r.templates[layoutName] = templ
	}

	return nil
}

func (r *Renderer) Render(w http.ResponseWriter, templateName string, data map[string]any) error {
	combinedData := maps.Clone(r.defaultData)
	templateName = templateName + r.templateSuffix

	maps.Copy(combinedData, data)

	t, ok := r.templates[templateName]
	if !ok {
		return fmt.Errorf("Template '%s' doesn't exist!", templateName)
	}

	return t.Execute(w, combinedData)
}

func (r *Renderer) handleError(w http.ResponseWriter, err error) {
	fmt.Fprintf(os.Stderr, "template error: %v\n", err)
	http.Error(w, "Internal Server Error", 500)
}
