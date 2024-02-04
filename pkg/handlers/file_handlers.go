package handlers

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
)

var chunkSize = 4096

type FileHandlers struct {
	saveToDir string

	logger *slog.Logger
}

func newFileHandlers(logger *slog.Logger, saveToDir string) (h *FileHandlers) {
	h = new(FileHandlers)
	h.saveToDir = saveToDir

	h.logger = logger

	return
}

var bufferPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:

		return bytes.NewBuffer(make([]byte, chunkSize))
	},
}

// Handle is responsible for handling incoming files requests
func (h *FileHandlers) Handle(w http.ResponseWriter, r *http.Request) {
	multiPartReader, err := r.MultipartReader()
	if err != nil {
		helpers.WrapBoth(w, err)
		return
	}

	dataBuffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		dataBuffer.Reset()
		bufferPool.Put(dataBuffer)
	}()

	err = h.handleMultiPartReader(multiPartReader, dataBuffer, h.saveToDir)
	if err != nil {
		helpers.WrapBoth(w, err)
	}
}

func (h *FileHandlers) handleMultiPartReader(
	multiPartReader *multipart.Reader,
	bytesBuffer *bytes.Buffer,
	saveDir string,
) error {
	var err error
	for {
		nextPart, err := multiPartReader.NextPart()
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

		err = h.writeDataToFile(bytesBuffer, multipartReadBuffer, fileWriteBuffer)
		if err != nil {
			h.logger.Error("error writing data to file", err)
			continue
		}
	}

	return err
}

// writeFlusher is an interface derived from io.Writer with one additional Flush method.
type writeFlusher interface {
	io.Writer
	Flush() error
}

func (h *FileHandlers) writeDataToFile(
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
