package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JoaoGabriel-Lima/desafio_devops/app1/internal/cache"
)

const DEFAULT_CACHE_TTL = 10 * time.Second

type Server struct {
	cache cache.Interface
}

func New(c cache.Interface) *Server {
	return &Server{
		cache: c,
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Desafio DevOps - João Gabriel Lima Marinho")
}

func (s *Server) handleStaticText(w http.ResponseWriter, r *http.Request) {
	const staticTextKey string = "texto_estatico"
	const ttl = DEFAULT_CACHE_TTL
	var valorCache, tempoParaExpirar, encontrado = s.cache.Get(staticTextKey)
	if encontrado {
		var tempoRestanteExpirar = time.Until(time.Unix(0, tempoParaExpirar))
		log.Printf("CACHE HIT (GO): texto estático encontrado no cache, tempo restante para expiração: %s", tempoRestanteExpirar)
		fmt.Fprintf(w, "Texto estático do cache: %s\n", valorCache)
		return
	}
	log.Println("CACHE MISS (GO): texto estático não encontrado no cache, usando valor padrão")
	s.cache.Set(staticTextKey, "Texto estático", ttl)

	fmt.Fprintln(w, "Texto estático (Go)")
}

func (s *Server) handleTime(w http.ResponseWriter, r *http.Request) {
	const cacheKey string = "hora_atual"
	const ttl = DEFAULT_CACHE_TTL
	var valorCache, tempoParaExpirar, encontrado = s.cache.Get(cacheKey)
	if encontrado {
		var tempoRestanteExpirar = time.Until(time.Unix(0, tempoParaExpirar))
		log.Printf("CACHE HIT (GO): hora atual encontrada no cache, tempo restante para expiração: %s", tempoRestanteExpirar)
		fmt.Fprintf(w, "Hora guardada na cache em memória: %s\n", valorCache)
		return
	}

	log.Println("CACHE MISS (GO): hora atual não encontrada no cache, gerando nova hora")
	var horaAtual string = time.Now().Format(time.RFC1123)
	s.cache.Set(cacheKey, horaAtual, ttl)
	response := fmt.Sprintf("Hora atual do servidor (GO): %s", horaAtual)
	fmt.Fprintln(w, response)
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static-text", s.handleStaticText)
	mux.HandleFunc("/time", s.handleTime)
}

