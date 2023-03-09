package multiplexer

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
	"github.com/i0tool5/simpleuploader/pkg/templates"
)

type Handlers struct {
	saveToDir string
	templates map[string]string
}

func NewHandlers(saveToDir string) (h *Handlers) {
	h = new(Handlers)
	h.saveToDir = saveToDir

	h.templates = make(map[string]string)

	return
}

// PrepareTemplates pre-creates template instances
func (h *Handlers) PrepareTemplates() (err error) {

	var template []byte

	template, err = templates.Temlpates.ReadFile("html/formtemplate.html")
	if err != nil {
		return
	}
	h.templates["form"] = string(template)

	template, err = templates.Fonts.ReadFile("fonts/Genos.ttf")
	if err != nil {
		return
	}
	h.templates["genos"] = string(template)

	return
}

func (h *Handlers) HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.templates["form"])
}

func (h *Handlers) HandleFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10_000)
	if err != nil {
		helpers.WrapBoth(w, err)
		return
	}

	files := r.MultipartForm.File
	if len(files["file"]) == 0 {
		log.Println("[!] No files to save.")
		w.WriteHeader(http.StatusBadRequest)
		helpers.WrapHttp(w)
		return
	}

	for _, file := range files["file"] {
		cf, err := file.Open()
		if err != nil {
			helpers.WrapBoth(w, err)
			return
		}
		defer cf.Close()
		fileContent, _ := io.ReadAll(cf)
		err = helpers.CreateFile(h.saveToDir, file.Filename, fileContent)
		if err != nil {
			log.Printf("[!] Error occured creating file %s\n", err)
		}
	}

	fmt.Fprint(w, h.templates["success"])
}

func (h *Handlers) HandleFonts(w http.ResponseWriter, r *http.Request) {
	font := h.templates["genos"]
	fmt.Fprint(w, font)
}
