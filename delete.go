package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// DeleteHandler handles deletes by objectId
func DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	if strings.HasPrefix(r.URL.Path, "/data") {
		requestPath := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if len(requestPath) < 2 || requestPath[1] == "" {
			err := New("Specify Object Id")
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		}

		var e entry
		e.ObjectId = bson.ObjectId(bson.ObjectIdHex(requestPath[1]))

		err := e.delete()
		if err != nil {
			http.Error(w, "While deleting: "+err.Error(), http.StatusNotFound)
			return err
		}

		data, err := load(e.Series)
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

	err := New("Specify Scope and Object Id")
	http.Error(w, err.Error(), http.StatusNotFound)
	return err
}
