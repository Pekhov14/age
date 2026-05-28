package main

import (
	"age/cmd"
	"age/internal/service"
	"age/internal/storage"
)

func main() {
	repo := storage.NewJsonRepository("data.json")
	serv := service.NewPersonService(repo)

	cmd.Execute(serv)
}
