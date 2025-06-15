# Desafio DevOps - João Gabriel Lima Marinho
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://www.python.org/)
[![Nginx](https://img.shields.io/badge/Nginx-009639?style=for-the-badge&logo=nginx&logoColor=white)](https://nginx.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)

Este repositório contém a implementação completa de um desafio DevOps com duas aplicações em linguagens diferentes, camadas de cache, observabilidade e infraestrutura automatizada.

## 📋 Objetivo do Desafio

Criar uma infraestrutura robusta com:

- Duas aplicações em linguagens diferentes (Go e Python)
- Camadas de cache com tempos de expiração distintos
- Observabilidade (Prometheus e Grafana)
- Facilidade de execução (Docker Compose)

## 🏗️ Arquitetura da Solução

<img src="./assets/infra_arq.png" alt="Arquitetura do Desafio DevOps" width="800"/>

## 🚀 Como Executar

### Pré-requisitos

- Docker
- Docker Compose

### Execução

```bash
# Clone o repositório
git clone <repository-url>
cd desafio_devops

# Iniciar toda a infraestrutura
docker-compose up -d

# Verificar status dos containers
docker-compose ps

# Visualizar logs
docker-compose logs -f
```

### Parar a aplicação

```bash
docker-compose down
```

## 🌐 Endpoints Disponíveis

### Aplicação Go (via Nginx)

- **GET** `http://localhost/go/` - Página inicial da aplicação Go
- **GET** `http://localhost/go/static-text` - Texto fixo (cache em memória: 10s)
- **GET** `http://localhost/go/time` - Horário atual do servidor (cache em memória: 10s)
- **GET** `http://localhost/go/metrics` - Métricas do Prometheus

### Aplicação Python (via Nginx)

- **GET** `http://localhost/python/` - Página inicial da aplicação Python
- **GET** `http://localhost/python/static-text` - Texto fixo (cache Nginx: 1min)
- **GET** `http://localhost/python/time` - Horário atual do servidor (cache Nginx: 1min)
- **GET** `http://localhost/python/metrics` - Métricas do Prometheus

### Observabilidade

- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)

## 🔧 Configuração das Camadas de Cache

### 1. Cache em Memória (Aplicação Go)

**Implementação**: Cache interno sincronizado com `sync.RWMutex`

**Configuração**:

```go
type CacheEmMemoria struct {
    sync.RWMutex
    items map[string]cacheItem
}

type cacheItem struct {
    valor      string
    expiracao  int64
}

const DEFAULT_CACHE_TTL = 10 * time.Second
```

**Características**:

- TTL: 10 segundos
- Thread-safe com RWMutex
- Expiração automática baseada em timestamp
- Cleanup automático no método Get()

### 2. Nginx Proxy Cache (Aplicação Python)

**Configuração no nginx.conf**:

```nginx
proxy_cache_path /var/cache/nginx keys_zone=my_cache:2m inactive=2m max_size=4m;

location /python/ {
    proxy_cache my_cache;
    proxy_cache_valid 200 1m;
    proxy_cache_key "$scheme$request_method$host$request_uri";
    add_header X-Cache-Status $upstream_cache_status;
    
    proxy_pass http://python-app:5000;
}
```

**Características**:

- TTL: 1 minuto para respostas 200
- Zona de cache: 2MB para chaves
- Cache máximo: 4MB
- Inatividade: 2 minutos
- Header `X-Cache-Status` para debug

## 📊 Observabilidade

### Métricas Coletadas

Ambas as aplicações expõem métricas no formato Prometheus:

**Contador de Requisições HTTP**:

```prometheus
http_total_requests{path="/", method="GET"}
```

### Dashboard Grafana

- **Go App**: Taxa de requisições por minuto
- **Python App**: Taxa de requisições por minuto
- Configuração automática via provisioning

## 🗂️ Estrutura do Projeto

```
desafio_devops/
├── docker-compose.yml              # Orquestração da infraestrutura
├── README.md                       # Este arquivo
├── app1_golang/                    # Aplicação Go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/app/main.go
│   └── internal/
│       ├── cache/cache.go          # Implementação do cache em memória
│       └── server/server.go        # Servidor HTTP com métricas
├── app2_python/                    # Aplicação Python
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app.py                      # Flask app com métricas
├── nginx/                          # Reverse proxy e cache
│   ├── Dockerfile
│   └── nginx.conf                  # Configuração com proxy cache
├── prometheus/                     # Monitoramento
│   └── prometheus.yml
└── grafana/                        # Visualização
    ├── dashboards/main-dashboard.json
    └── provisioning/
        ├── dashboards/provider.yml
        └── datasources/prometheus.yml
```

## 🔄 Fluxo de Atualização de Componentes


## 🚀 Pontos de Melhoria

### Melhorias do Aplicativo Go

### Melhorias do Aplicativo Python

### Melhorias Gerais da Infraestrutura


## 🧪 Testes

```bash
# Testar cache do Go (10 segundos)
curl http://localhost/go/time
curl http://localhost/go/time  # Deve retornar valor do cache

# Testar cache do Nginx (1 minuto)
curl -H "X-Cache-Status: debug" http://localhost/python/time
```

## 📝 Notas Técnicas

- **Go**: Utiliza goroutines para concorrência e RWMutex para thread-safety
- **Python**: Flask com Werkzeug middleware para métricas
- **Nginx**: Configurado como reverse proxy com cache layer
- **Prometheus**: Scraping automático das métricas das aplicações
- **Grafana**: Dashboards provisionados automaticamente

<!-- ## 🤝 Contribuições

Este projeto foi desenvolvido como parte de um desafio DevOps, demonstrando conhecimentos em:

- Containerização com Docker
- Orquestração com Docker Compose
- Desenvolvimento em Go e Python
- Estratégias de cache em múltiplas camadas
- Observabilidade com Prometheus/Grafana
- Reverse proxy com Nginx -->
