package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hararudoka/shrt/model"
	"github.com/hararudoka/shrt/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	*service.Service
	*zerolog.Logger
}

func New(s service.Service) http.Handler {
	return Handler{
		&s,
		&log.Logger,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.Path == "/api/short" {
			h.Short(w, r)
			return
		}
		if r.URL.Path == "/api/url" {
			h.URL(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./view/index.html")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/assets") {
			http.FileServer(http.Dir("./view")).ServeHTTP(w, r)
			return
		}

		if len(r.URL.Path) == 7 {
			h.Redirect(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

// Short returns JSON with a short. Asking for a URL in the body.
func (h *Handler) Short(w http.ResponseWriter, r *http.Request) {
	var u model.URL
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Info().Err(err).Msg("error during json decoding")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the HTTP request with the request body using ZeroLog
	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Str("URL", u.URL). // Log the request body as a string
		Msg("Received HTTP request")

	short, err := h.Service.URL2Hash(u.URL)
	if err != nil {
		log.Info().Err(err).Msg("error during h.Service.URL2Hash")
		// we don't want to show internal errors to the user, so we return 404
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(model.Short{Short: short})
	if err != nil {
		log.Info().Err(err).Msg("error during json decoding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// URL returns JSON with a full URL. Asking for a Short in the body.
func (h *Handler) URL(w http.ResponseWriter, r *http.Request) {
	var s model.Short
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Info().Err(err).Msg("error during Body json decoding")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the HTTP request with the request body using ZeroLog
	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Str("URL", s.Short). // Log the request body as a string
		Msg("Received HTTP request")

	url, err := h.Service.Hash2URL(s.Short)
	if err != nil {
		log.Info().Err(err).Msg("error during h.Service.Hash2URL")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(model.URL{URL: url})
	if err != nil {
		log.Info().Err(err).Msg("error during URL json decoding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// we here only if it is GET method of "/"
	url, err := h.Service.Hash2URL(r.URL.Path[1:])
	if err != nil {
		log.Info().Err(err).Msg("redirect error")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "https://"+url, http.StatusFound)
}
