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
		personInfos = append(personInfos, buildPersonInfo(person.Name, person.Birthday))
	}

	return personInfos, nil
}

func (s *PersonService) Delete(name string) error {
	return s.repo.DeleteByName(name)
}

func (s *PersonService) Update(oldName, newName string, birthday time.Time) error {
	person := storage.Person{Name: newName, Birthday: birthday}
	return s.repo.Update(oldName, person)
}

func (s *PersonService) PreviewBirthday(birthday time.Time) PersonInfo {
	return buildPersonInfo("Preview", birthday)
}

func buildPersonInfo(name string, birthday time.Time) PersonInfo {
	strFormat, days := calculateDaysUntilBD(birthday)

	return PersonInfo{
		Name:        name,
		Birthday:    birthday.Format("2006-01-02"),
		Age:         calculateAge(birthday),
		DaysToBirth: days,
		DaysUntilBD: strFormat,
	}
}
