package nethttp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hararudoka/shrt/internal/model"
	"github.com/hararudoka/shrt/internal/service"
)

type Handler struct {
	*service.Service
}

func New(s service.Service) http.Handler {
	return Handler{
		&s,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if len(r.URL.Path) == 7 {
			h.Redirect(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		if r.URL.Path == "/short" {
			h.Short(w, r)
			return
		}
		if r.URL.Path == "/url" {
			h.URL(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// Short returns JSON with a short. Asking for a URL in the body.
func (h Handler) Short(w http.ResponseWriter, r *http.Request) {
	var u model.URL
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	short, err := h.Service.URL2Hash(u.URL)
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(model.Short{Short: short})
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// URL returns JSON with a full URL. Asking for a Short in the body.
func (h Handler) URL(w http.ResponseWriter, r *http.Request) {
	var s model.Short
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := h.Service.Hash2URL(s.Short)
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(model.URL{URL: url})
	if err != nil {
		// TODO: errors
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// we here only if it is GET method of "/"

	url, err := h.Service.Hash2URL(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "https://"+url, http.StatusFound)
}
