package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JoaoGabriel-Lima/desafio_devops/app1/internal/cache"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Tempo padrão de expiração do cache
const DEFAULT_CACHE_TTL = 10 * time.Second


// httpTotalRequets é um contador Prometheus para rastrear o número total de requisições HTTP recebidas
var httpTotalRequets = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_total_requests",
		Help: "Número total de requisições HTTP recebidas",
	},
	[]string{"path", "method"},
)


func init() {
	prometheus.MustRegister(httpTotalRequets)
}

// Server representa o servidor HTTP que lida com as rotas e interage com o cache
type Server struct {
	cache cache.Interface
}

// NewServer cria uma nova instância do servidor HTTP com o cache fornecido
func New(c cache.Interface) *Server {
	return &Server{
		cache: c,
	}
}

// handleIndex é o manipulador para a rota raiz ("/") que retorna uma mensagem de boas-vindas
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	httpTotalRequets.WithLabelValues(r.URL.Path, r.Method).Inc()
	
	fmt.Fprintln(w, "Desafio DevOps - João Gabriel Lima Marinho - Servidor GO")
}

// handleStaticText é o manipulador para a rota "/static-text" que retorna um texto estático
func (s *Server) handleStaticText(w http.ResponseWriter, r *http.Request) {
	httpTotalRequets.WithLabelValues(r.URL.Path, r.Method).Inc()

	// Faço uma verificação no cache para ver se o texto estático já foi armazenado, se sim, retorna o valor do cache, se não, armazena o valor no cache e retorna o valor
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

	fmt.Fprintln(w, "Texto estático (GO)")
}

// handleTime é o manipulador para a rota "/time" que retorna a hora atual do servidor, armazenando em cache
func (s *Server) handleTime(w http.ResponseWriter, r *http.Request) {
	httpTotalRequets.WithLabelValues(r.URL.Path, r.Method).Inc()

	// Faço uma verificação no cache para ver se a hora atual já foi armazenada, se sim, retorna o valor do cache, se não, armazena a hora atual no cache e retorna o valor
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

// RegisterRoutes registra as rotas do servidor HTTP
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static-text", s.handleStaticText)
	mux.HandleFunc("/time", s.handleTime)
	mux.Handle("/metrics", promhttp.Handler())
}

