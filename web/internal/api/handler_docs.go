package api

import (
	"html/template"
	"net/http"

	"github.com/kiselevos/new_tax/pkg/logx"
)

type EndpointInfo struct {
	Method             string
	Path               string
	Description        string
	ExampleRequestFull string
	ExampleRequestLite string
	ExampleResponse    string
	Group              string
}

type DocsHandler struct {
	Endpoints []EndpointInfo
	ApiVers   string
	Template  *template.Template
}

func NewDocsHandler(apiVers string, tmpl *template.Template) *DocsHandler {
	return &DocsHandler{
		Endpoints: []EndpointInfo{},
		ApiVers:   apiVers,
		Template:  tmpl,
	}
}

func (h *DocsHandler) HandleApiDocs(w http.ResponseWriter, r *http.Request) {
	data := struct {
		APIVersion string
		Endpoints  []EndpointInfo
	}{
		APIVersion: h.ApiVers,
		Endpoints:  h.Endpoints,
	}

	if err := h.Template.ExecuteTemplate(w, "api", data); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (h *DocsHandler) AddEndpoint(e EndpointInfo) {
	h.Endpoints = append(h.Endpoints, e)
}
