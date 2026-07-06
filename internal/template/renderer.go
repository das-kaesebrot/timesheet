package template

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/utility"
)

type Renderer struct {
	dir       string
	funcs     template.FuncMap
	layouts   *template.Template
	templates map[string]string
	weekdays  []time.Weekday
	timezones []string
	version   string
}

func New(dir string, version string) (*Renderer, error) {
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
		},
		weekdays: utility.GetWeekdays(),
	}

	availableTimezones, err := utility.GetAllTimezones(true)
	if err != nil {
		return nil, err
	}
	r.timezones = availableTimezones

	log.Printf("Found OS timezones: %v", availableTimezones)

	r.version = version

	return r, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	partials, err := filepath.Glob(filepath.Join(r.dir, "partials", "*.html"))
	if err != nil {
		r.handleError(w, err)
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		dataMap["Version"] = r.version
		dataMap["Timezones"] = r.timezones
		dataMap["Weekdays"] = r.weekdays
	}

	// ugly but this way we keep the strict order
	templateFiles := append(append([]string{filepath.Join(r.dir, "layouts", "base.html")}, partials...), filepath.Join(r.dir, name+".html"))

	t := template.New("base.html").Funcs(r.funcs)
	t, err = t.ParseFiles(templateFiles...)
	if err != nil {
		r.handleError(w, err)
	}

	err = t.Execute(w, data)

	if err != nil {
		r.handleError(w, err)
	}
}

func (r *Renderer) handleError(w http.ResponseWriter, err error) {
	fmt.Fprintf(os.Stderr, "template error: %v\n", err)
	http.Error(w, "Internal Server Error", 500)
}
