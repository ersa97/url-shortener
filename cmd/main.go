package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

type Link struct {
	Target string `json:"target"`
	Short  string `json:"Short"`
}

func (l *Link) Bind(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&l)
}

func (l *Link) Validate() error {
	if l.Target == "" {
		return errors.New("target cannot be empty")
	}
	return nil
}

type Response struct {
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

var LinkDB = []Link{}

func handleResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Message: message,
		Data:    data,
	})
}

func methodHelper(r *http.Request, method string) bool {
	return r.Method == method
}

func main() {
	http.HandleFunc("/add", AddLink)
	http.HandleFunc("/fetch", FetchLink)

	http.ListenAndServe(":3030", nil)
}

func AddLink(w http.ResponseWriter, r *http.Request) {
	ok := methodHelper(r, "POST")
	if !ok {
		handleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	var newLink Link
	if err := newLink.Bind(r); err != nil {
		handleResponse(w, http.StatusInternalServerError, "cannot get body", nil)
		return
	}
	if err := newLink.Validate(); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	for _, v := range LinkDB {
		if v.Target == newLink.Target {
			handleResponse(w, http.StatusBadRequest, "link already stored", nil)
			return
		}
	}

	genString := fmt.Sprint(rand.Int63n(1000))
	newLink.Short = fmt.Sprintf("bit.ly/%s", genString)

	LinkDB = append(LinkDB, newLink)

	handleResponse(w, http.StatusOK, "short link generated", newLink)
}

func FetchLink(w http.ResponseWriter, r *http.Request) {
	ok := methodHelper(r, "GET")
	if !ok {
		handleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	short := r.URL.Query().Get("short")
	var res *Link
	for _, v := range LinkDB {
		if v.Short == short {
			res = &Link{
				Target: v.Target,
				Short:  short,
			}
		}
	}
	if res == nil {
		handleResponse(w, http.StatusNotFound, "link not found", nil)
		return
	}

	handleResponse(w, http.StatusOK, "link found", res)
}
