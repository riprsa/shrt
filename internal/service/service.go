package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/hararudoka/shrt/internal/model"
)

type Storage interface {
	Insert(URL, short string) error
	ByShort(short string) (model.Data, error)
	ByURL(url string) (model.Data, error)
}

// Service is a struct for service layer
type Service struct {
	db Storage
}

// New creates new Service
func New(db Storage) *Service {
	return &Service{
		db: db,
	}
}

// URL2Hash takes a full URL, sanitizes it and returns a short URL hash.
// A new URL hash is created and stored, while for an existing URL
// such hash is looked up from the storage.
func (s Service) URL2Hash(url string) (string, error) {
	// validate url
	url, ok := s.ValidateURL(url)
	if !ok {
		return "", fmt.Errorf("url is broken")
	}

	// check if url exists in DB
	dataFromDB, err := s.db.ByURL(url)
	if err == sql.ErrNoRows {
		short, err := s.createShort(url)
		if err != nil {
			return "", err
		}
		return short, nil
	} else if err != nil {
		return "", err
	}

	return dataFromDB.Short, nil
}

// Hash2URL returns full URS by the short link.
func (s Service) Hash2URL(short string) (string, error) {
	data, err := s.db.ByShort(short)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("404 page not found")
	} else if err != nil {
		return "", err
	}

	url, ok := s.ValidateURL(data.URL)
	if !ok {
		return "", fmt.Errorf("url is broken")
	}
	return url, nil
}

func (s Service) createShort(url string) (string, error) {
	for {
		short := generateShort()
		if yes, err := s.exist(short); !yes {
			if err != nil {
				return "", err
			}
			err = s.db.Insert(url, short)
			if err != nil {
				return "", err
			}
			return short, nil
		}
	}
}

// exist checks if short url exists in DB
func (s Service) exist(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}

// generateShort generates short url
func generateShort() string {
	var s []byte
	ra := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 6; i++ {
		s = append(s, ra[rand.Intn(len(ra))])
	}
	return string(s)
}

var ErrEmptyURL = errors.New("URL must not be empty string")

// SanitizeURL return URL string with properly encoded query and stripped of
// schema and user identification parts. Sanitized URL contains only host,
// path, query and fragment parts (per RFC 3986) in return to the provided
// ID string.
func SanitizeURL(u string) (string, error) {
	if u == "" {
		return "", ErrEmptyURL
	}
	url, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	// Strip URL from the User identification and from the scheme
	// prefix before returning.
	url.User = nil
	return strings.SplitN(url.String(), "://", 2)[1], nil
}

// ValidateURL validates and adds protocol if link hasn't it
func (s Service) ValidateURL(url string) (string, bool) {
	// ok := s.URLRegexp.Match([]byte(url))
	// if !ok {
	// 	return "", false
	// }

	// withProtocol := s.ProtocolRegexp.Match([]byte(url))
	// if !withProtocol {
	// 	url = "http://" + url
	// }

	return url, true
}
