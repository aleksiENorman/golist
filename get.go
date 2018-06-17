package main

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Get handler transforms csv data into json and send it via HTTP socket
func GetHandler(w http.ResponseWriter, r *http.Request) error {

	if strings.HasPrefix(r.URL.Path, "/data") {
		w.Header().Set("Content-Type", "application/json")

		requestPath := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		var series string

		if len(requestPath) < 2 || requestPath[1] == "" {
			series = "index"
		} else {
			series = requestPath[1]
		}

		data, err := load(series)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return nil
	}

	return StaticHandler(w, r)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) error {
	mimetype := mime.TypeByExtension(filepath.Ext(getStaticFilename(r)))
	w.Header().Set("Content-Type", mimetype)

	handle, err := os.Open(getStaticFilename(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, handle)
	return err
}

func getStaticFilename(r *http.Request) string {
	if r.URL.Path == "/" {
		return "static/index.html"
	}

	return "static" + r.URL.Path
}
