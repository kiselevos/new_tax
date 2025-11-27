package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kiselevos/new_tax/web"
)

func TestTemplatesRender(t *testing.T) {
	tmpls := loadTemplates(t)

	data := PrepareIndexData()

	if err := tmpls.ExecuteTemplate(httptest.NewRecorder(), "index", data); err != nil {
		t.Errorf("index template execution failed: %v", err)
	}

	for _, name := range []string{"about", "regional_info", "special_tax_modes"} {
		if err := tmpls.ExecuteTemplate(httptest.NewRecorder(), name, nil); err != nil {
			t.Errorf("%s template execution failed: %v", name, err)
		}
	}
}

func loadTemplates(t *testing.T) *template.Template {
	t.Helper()

	pattern := "templates/*.tmpl"
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		pattern = filepath.Join("..", "templates", "*.tmpl")
	}

	tmpls, err := template.New("").Funcs(web.Funcs).ParseGlob(pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}
	return tmpls
}

func TestHandlers_StatusOK(t *testing.T) {
	s := &Server{Tmpl: loadTemplates(t)}

	tests := []struct {
		name string
		path string
		fn   func(http.ResponseWriter, *http.Request)
	}{
		{"index", "/", s.Index},
		{"about", "/about", s.About},
		{"regional_info", "/regional-info", s.RegionalInfo},
		{"special_tax_modes", "/special-tax-modes", s.SpecialTaxModes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			tt.fn(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("%s handler returned status %v", tt.name, w.Code)
			}
		})
	}
}
