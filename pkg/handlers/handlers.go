package handlers

import (
	"log/slog"

	"github.com/i0tool5/simpleuploader/pkg/templates"
)

type Handlers struct {
	Static *StaticHandlers
	File   *FileHandlers
}

func New(
	saveDir string,
	logger *slog.Logger,
	template *templates.Template,
) Handlers {
	return Handlers{
		Static: newStaticHandlers(template, logger),
		File:   newFileHandlers(logger, saveDir),
	}
}
