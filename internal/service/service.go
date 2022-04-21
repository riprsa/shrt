package service

import (
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"

	"shorter/internal/storage"
)

// Service is a struct for service layer
type Service struct {
	db             *storage.DB
	URLRegexp      *regexp.Regexp
	ProtocolRegexp *regexp.Regexp
}

// New creates new Service
func New(db *storage.DB) *Service {
	return &Service{
		db:             db,
		URLRegexp:      regexp.MustCompile("([a-z]*:\\/\\/)?[a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)"),
		ProtocolRegexp: regexp.MustCompile("[a-z]*:\\/\\/[a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)"),
	}
}

// GetShort generates or takes from DB Short
func (s Service) GetShort(url string) (string, error) {
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

func (s Service) GetURL(short string) (string, error) {
	data, err := s.db.ByShort(short)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("not found")
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

// ValidateURL validates and adds protocol if link hasn't it
func (s Service) ValidateURL(url string) (string, bool) {
	ok := s.URLRegexp.Match([]byte(url))
	if !ok {
		return "", false
	}

	withProtocol := s.ProtocolRegexp.Match([]byte(url))
	if !withProtocol {
		url = "http://" + url
	}

	return url, true
}
