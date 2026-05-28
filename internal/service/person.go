package service

import (
	"age/internal/storage"
	"time"
)

type PersonInfo struct {
	Name        string
	Birthday    string
	Age         int
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
		personInfos = append(personInfos, PersonInfo{
			Name:        person.Name,
			Birthday:    person.Birthday.Format("2006-01-02"),
			Age:         calculateAge(person.Birthday),
			DaysUntilBD: calculateDaysUnliBd(person.Birthday),
		})
	}

	return personInfos, nil
}

func calculateDaysUnliBd(birthday time.Time) string {
	// todo: create calculateDaysUnliBd
	return "Next week"
}

func calculateAge(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	// If your birthday hasn't arrived yet this year, subtract 1.
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}
