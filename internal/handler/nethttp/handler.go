package nethttp

import (
	"net/http"

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

func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if len(r.URL.Path) == 7 {
			h.GetRedirect(rw, r)
		}
	}
	if r.Method == http.MethodPost {
		if r.URL.Path == "short" {
			h.GetURL(rw, r)
		}
		if r.URL.Path == "url" {
			h.GetShort(rw, r)
		}
	}
}

// GetShort returns JSON with a short. Asking for a URL in the body.
func (h Handler) GetShort(rw http.ResponseWriter, r *http.Request) {
	// url := r
}

// GetURL returns JSON with a full URL. Asking for a Short in the body.
func (h Handler) GetURL(rw http.ResponseWriter, r *http.Request) {
	// short := r.URL.Path[1:]
}

//
func (h Handler) GetRedirect(rw http.ResponseWriter, r *http.Request) {
	// we here only if it is GET method of "/"

	url, err := h.Service.Hash2URL(r.URL.Path[1:])
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	http.Redirect(rw, r, url, http.StatusFound)
}
