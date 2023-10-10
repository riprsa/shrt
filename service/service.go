package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	Insert(URL, short string) error
	ByShort(short string) (string, error)
	ByURL(url string) (string, error)
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
	url, err := SanitizeURL(url)
	if err != nil {
		return "", err
	}

	// check if url exists in DB
	short, err := s.db.ByURL(url)
	if err == pgx.ErrNoRows {
		return s.createShort(url)
	}
	return short, err
}

// Hash2URL returns full URL by the short link.
func (s Service) Hash2URL(short string) (string, error) {
	url, err := s.db.ByShort(short)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("404 page not found")
	} else if err != nil {
		return "", err
	}

	url, err = SanitizeURL(url)
	if err != nil {
		return "", err
	}
	return url, nil
}

// createShort creates short URL and stores it in DB
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
func (s Service) exist(short string) (bool, error) {
	_, err := s.db.ByShort(short)
	if err == pgx.ErrNoRows {
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
func SanitizeURL(u string) (string, error) { // TODO: check to ensure that is works, expand test cases
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
	url.Scheme = "http"

	splitted := strings.Split(url.String(), "://")
	if len(splitted) != 2 {
		return "", fmt.Errorf("invalid URL")
	}

	return splitted[1], nil
}
