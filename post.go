package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

// PostHandler converts post request body to string
func PostHandler(w http.ResponseWriter, r *http.Request) entry {
	var result entry

	if r.Method == "POST" {
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if strings.HasPrefix(mediaType, "multipart/") {
			result, err = multipartRead(r, params["boundary"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return result
			}
		} else if mediaType == "application/x-www-form-urlencoded" {
			result, err = urlEncodedRead(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return result
			}
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return result
	}

	if err := result.saveValidate(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return result
	}

	if err := result.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return result
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return result
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)

	return result
}

// Read entry from "application/x-www-form-urlencoded"
// Post request
func urlEncodedRead(r *http.Request) (entry, error) {
	var result entry

	if err := r.ParseForm(); err != nil {
		return result, err
	}

	for key, value := range r.Form {
		result.setIndex(key, []byte(value[0]))
	}

	return result, nil
}

// Read entry from "form/multipart"
// Post request
func multipartRead(r *http.Request, boundary string) (entry, error) {
	var result entry

	reader := multipart.NewReader(r.Body, boundary)
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}

		buffer, err := ioutil.ReadAll(p)
		if err != nil {
			return result, err
		}

		result.setIndex(p.FormName(), buffer)
	}

	return result, nil
}
