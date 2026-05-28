package service

import (
	"age/internal/storage"
	"time"
)

type PersonInfo struct {
	Name        string
	Birthday    string
	Age         int
	DaysToBirth int
	DaysUntilBD string
}

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

func (s *PersonService) List() ([]PersonInfo, error) {
	persons, err := s.repo.ListPeople()
	if err != nil {
		return nil, err
	}

	var personInfos []PersonInfo

	for _, person := range persons {
		strFormat, days := calculateDaysUntilBD(person.Birthday)

		personInfos = append(personInfos, PersonInfo{
			Name:        person.Name,
			Birthday:    person.Birthday.Format("2006-01-02"),
			Age:         calculateAge(person.Birthday),
			DaysToBirth: days,
			DaysUntilBD: strFormat,
		})
	}

	return personInfos, nil
}
