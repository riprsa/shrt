package service

import (
	"database/sql"
	"fmt"
	"math/rand"

	"shorter/internal/storage"
	"shorter/internal/validate"
)

type Service struct {
	db *storage.DB
}

func New(db *storage.DB) *Service {
	return &Service{db: db}
}

// GetShort generates or take from DB Short
func (s Service) GetShort(url string) (string, error)  {
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

func (s Service) GetURL(short string) (string, error)  {
	data, err := s.db.ByShort(short)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("not found: ")
	} else if err != nil {
		return "", err
	}

	url, ok := validate.URL(data.URL)
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

func (s Service) exist(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}

func generateShort() string {
	var s []byte
	ra := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 6; i++ {
		s = append(s, ra[rand.Intn(len(ra))])
	}
	return string(s)
}
