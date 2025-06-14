package main

import (
	"log"
	"net/http"

	"github.com/JoaoGabriel-Lima/desafio_devops/app1/internal/cache"
	"github.com/JoaoGabriel-Lima/desafio_devops/app1/internal/server"
)

const SERVER_PORT string = ":8080"

func main() {
	var mux *http.ServeMux = http.NewServeMux()	
	var cache cache.Interface = cache.NewCache()
	
	var server = server.New(cache)
	server.RegisterRoutes(mux)
	
	// Iniciando o servidor HTTP
	log.Printf("Servidor iniciado na porta %s", SERVER_PORT)
	var err = http.ListenAndServe(SERVER_PORT, mux);
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}