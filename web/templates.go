package web

import (
	"encoding/json"
	"net/http"
	"text/template"
)

type TemplatesHandler struct {
	t map[string]*template.Template
}

func NewTemplatesHandler(templates map[string]*template.Template) *TemplatesHandler {
	return &TemplatesHandler{t: templates}
}

func (h *TemplatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(h.t))
	for s, _ := range h.t {
		keys = append(keys, s)
	}
	bytes, _ := json.Marshal(keys)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
