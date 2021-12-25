package middleware

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type HTTPUploads struct {
	dir string
	mem int64
}

func NewHTTPUploads(dir string, memLimit int64) *HTTPUploads {
	h := &HTTPUploads{dir, memLimit}

	return h
}

func (m *HTTPUploads) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(m.mem)
	if r.MultipartForm == nil || r.MultipartForm.File["upload[]"] == nil {
		http.Error(w, "no files found", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["upload[]"]
	for _, fileHeader := range files {
		_, params, err := mime.ParseMediaType(fileHeader.Header.Get("Content-Disposition"))
		if err != nil {
			renderError(w, err, "error opening data", http.StatusBadRequest)
			return
		}

		dfile := filepath.Join(m.dir, params["filename"])
		ddir := filepath.Dir(dfile)

		file, err := fileHeader.Open()
		if err != nil {
			renderError(w, err, "error opening data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			renderError(w, err, "error reading data", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			renderError(w, err, "error reading data", http.StatusBadRequest)
			return
		}

		err = os.MkdirAll(ddir, os.ModePerm)
		if err != nil {
			renderError(w, err, "internal server error", http.StatusInternalServerError)
			return
		}

		f, err := os.Create(dfile)
		if err != nil {
			renderError(w, err, "internal server error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			renderError(w, err, "internal server error", http.StatusBadRequest)
			return
		}
	}
}

func renderError(w http.ResponseWriter, err error, msg string, code int) {
	http.Error(w, fmt.Sprintf("%s: %s", msg, err), code)
}
