package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/i0tool5/simpleuploader/pkg/templates"
)

type StaticHandlers struct {
	templates *templates.Template
	logger    *slog.Logger
}

func newStaticHandlers(t *templates.Template, l *slog.Logger) *StaticHandlers {
	return &StaticHandlers{
		templates: t,
		logger:    l,
	}
}

// HandleMain is responsible for handling requests to the main page.
func (h *StaticHandlers) HandleMain(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, h.templates.Form())
	if err != nil {
		h.logger.Error("can't handle main request", err)
	}
}

// HandleFonts is responsible for handling font requests.
func (h *StaticHandlers) HandleFonts(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, h.templates.Genos())
	if err != nil {
		h.logger.Error("can't send font template", err)
	}
}
