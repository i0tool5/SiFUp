package multiplexer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
	"github.com/i0tool5/simpleuploader/pkg/templates"
)

var chunkSize = 4096

type Handlers struct {
	saveToDir string
	templates map[string]string

	logger *slog.Logger
}

func NewHandlers(logger *slog.Logger, saveToDir string) (h *Handlers) {
	h = new(Handlers)
	h.saveToDir = saveToDir

	h.templates = make(map[string]string)
	h.logger = logger

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

// HandleFonts is responsible for handling font requests
func (h *Handlers) HandleFonts(w http.ResponseWriter, r *http.Request) {
	font := h.templates["genos"]
	fmt.Fprint(w, font)
}

// HandleFiles is responsible for handling incoming files requests
func (h *Handlers) HandleFiles(w http.ResponseWriter, r *http.Request) {
	multiPartReader, err := r.MultipartReader()
	if err != nil {
		helpers.WrapBoth(w, err)
		return
	}

	dataBuffer := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(dataBuffer)
	defer dataBuffer.Reset()

	err = h.handleMultiPartReader(multiPartReader, dataBuffer, h.saveToDir)
	if err != nil {
		helpers.WrapBoth(w, err)
	}

	fmt.Fprint(w, h.templates["success"])
}

func (h *Handlers) handleMultiPartReader(
	mpReader *multipart.Reader,
	buff *bytes.Buffer,
	saveDir string,
) error {
	var err error
	for {
		nextPart, err := mpReader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		fullpath := saveDir + nextPart.FileName()
		h.logger.Debug("creating file", slog.Any("path", fullpath))
		fileForWrite, err := os.Create(fullpath)
		if err != nil {
			return err
		}
		defer fileForWrite.Close()

		multipartReadBuffer := bufio.NewReader(nextPart)
		fileWriteBuffer := bufio.NewWriter(fileForWrite)
		err = h.writeDataToFile(buff, multipartReadBuffer, fileWriteBuffer)
		if err != nil {
			h.logger.Error("error writing data to file", err)
			continue
		}
	}

	return err
}

type writeFlusher interface {
	io.Writer
	Flush() error
}

func (h *Handlers) writeDataToFile(
	buffer *bytes.Buffer,
	dataReader io.Reader,
	dataWriter writeFlusher,
) error {
	for {
		dat := buffer.Bytes()
		n, err := dataReader.Read(dat)
		if err != nil && n == 0 {
			if errors.Is(err, io.EOF) {
				dataWriter.Flush()
				h.logger.Debug("got EOF reading data")
				break
			}
			return err
		}

		h.logger.Debug("writing data", slog.Any("data part", dat[:n]))
		_, err = dataWriter.Write(dat[:n])
		if err != nil {
			return err
		}
		dataWriter.Flush()
	}

	return nil
}
