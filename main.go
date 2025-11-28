package main

import (
	"log"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/routes"
)

func main() {
	if err := database.ConectaComBancoDeDados(); err != nil {
		log.Panicf("Erro ao iniciar aplicacao: %v", err)
	}
	routes.HandleRequest()
}
