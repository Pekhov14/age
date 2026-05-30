package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type JsonRepository struct {
	filePath string
}

func NewJsonRepository(filePath string) *JsonRepository {
	return &JsonRepository{filePath: filePath}
}

func (r *JsonRepository) AddPerson(p Person) error {
	persons, err := r.load()
	if err != nil {
		return err
	}

	persons = append(persons, p)

	return r.save(persons)
}

func (r *JsonRepository) DeleteByName(name string) error {
	persons, err := r.load()
	if err != nil {
		return err
	}

	lengthBeforeDelete := len(persons)

	persons = slices.DeleteFunc(persons, func(p Person) bool {
		return p.Name == name
	})

	lengthAfterDelete := len(persons)

	if lengthAfterDelete == lengthBeforeDelete {
		return fmt.Errorf("person %q not found", name)
	}

	return r.save(persons)
}

func (r *JsonRepository) ListPeople() ([]Person, error) {
	return r.load()
}

func (r *JsonRepository) load() ([]Person, error) {
	var persons []Person

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return []Person{}, nil // File not exists
	}

	err = json.Unmarshal(data, &persons)
	return persons, err
}

func (r *JsonRepository) save(people []Person) error {
	data, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}
