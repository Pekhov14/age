package service

import (
	"age/internal/storage"
	"time"
)

type PersonService struct {
	repo storage.Repository
}

func NewPersonService(repo storage.Repository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) Add(name string, birthday time.Time) error {
	p := storage.Person{Name: name, Birthday: birthday}
	return s.repo.AddPerson(p)
}
