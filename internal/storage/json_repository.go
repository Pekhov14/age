package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonRepository struct {
	filePath string
}

func NewJsonRepository(filePath string) *JsonRepository {
	return &JsonRepository{filePath: filePath}
}

func (r *JsonRepository) AddPerson(p Person) error {
	// people, _ := r.ListPeople()
	// people = append(people, p)
	// return r.save(people)

	//path := r.filePath

	// TODO: фильтровать можно в цикле есть ли такое имя
	// Или добавить еще ключ имя

	var persons []Person

	fileName := "data.json"

	data, err := os.ReadFile(r.filePath)
	if err == nil {
		_ = json.Unmarshal(data, &persons)
	}

	persons = append(persons, p)

	updatedData, _ := json.MarshalIndent(persons, "", "  ")
	_ = os.WriteFile(fileName, updatedData, 0644)

	fmt.Println("Типа добавил в файл)")
	return nil
}

func (r *JsonRepository) DeletePerson(p Person) error {
	// реализация позже
	return nil
}

func (r *JsonRepository) ListPeople() ([]Person, error) {
	// data, err := os.ReadFile(r.filePath)
	// if err != nil {
	//     return []Person{}, nil // File not exits
	// }
	// var people []Person
	// err = json.Unmarshal(data, &people)
	// return people, err

	return []Person{}, nil
}

func (r *JsonRepository) save(people []Person) error {
	// data, err := json.Marshal(people)
	// if err != nil {
	//     return err
	// }
	// return os.WriteFile(r.filePath, data, 0644)

	return nil
}
