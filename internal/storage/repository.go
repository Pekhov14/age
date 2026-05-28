package storage

import "time"

type Person struct {
	Name     string    `json:"name"`
	Birthday time.Time `json:"age"`
}

type Repository interface {
	AddPerson(Person) error
	DeletePerson(Person) error
	ListPeople() ([]Person, error)
}
