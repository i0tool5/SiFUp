package multiplexer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
	"github.com/i0tool5/simpleuploader/pkg/templates"
)

var chunkSize = 512

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

var bufPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:

		return bytes.NewBuffer(make([]byte, chunkSize))
	},
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

// HandleMain is responsible for handling requests to the main page.
func (h *Handlers) HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.templates["form"])
}

// HandleFiles is responsible for handling incoming files requests
func (h *Handlers) HandleFiles(w http.ResponseWriter, r *http.Request) {
	multiPartReader, err := r.MultipartReader()
	if err != nil {
		helpers.WrapBoth(w, err)
		return
	}

	buff := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buff)
	defer buff.Reset()

	err = handleMultiPartReader(multiPartReader, buff, h.saveToDir)
	if err != nil {
		helpers.WrapBoth(w, err)
	}

	fmt.Fprint(w, h.templates["success"])
}

// HandleFonts is responsible for handling font requests
func (h *Handlers) HandleFonts(w http.ResponseWriter, r *http.Request) {
	font := h.templates["genos"]
	fmt.Fprint(w, font)
}

func handleMultiPartReader(mpReader *multipart.Reader, buff *bytes.Buffer, saveDir string) error {
	for {
		nextPart, err := mpReader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		wfile, err := os.Create(saveDir + nextPart.FileName())
		if err != nil {
			return err
		}
		defer wfile.Close()

		bufRead := bufio.NewReader(nextPart)
		bufWriter := bufio.NewWriter(wfile)

		for {
			dat := buff.Bytes()
			n, err := bufRead.Read(dat)
			if err != nil && n == 0 {
				if errors.Is(err, io.EOF) {
					bufWriter.Flush()
					break
				}
				return err
			}

			_, err = bufWriter.Write(dat[:n])
			if err != nil {
				return err
			}
			bufWriter.Flush()
		}
	}

	return nil
}
