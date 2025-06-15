# Desafio DevOps - João Gabriel Lima Marinho
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://www.python.org/)
[![Nginx](https://img.shields.io/badge/Nginx-009639?style=for-the-badge&logo=nginx&logoColor=white)](https://nginx.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)

Este repositório contém a implementação completa de um desafio DevOps com duas aplicações em linguagens diferentes, camadas de cache, observabilidade e infraestrutura automatizada.

## 🔀 Implementações Disponíveis

### Implementação Principal (Branch Main)

- **Go**: Cache em memória interno
- **Python**: Nginx Proxy Cache

### Implementação Alternativa (Branch Redis)

- **Go**: Cache em memória interno
- **Python**: Redis Cache

🔗 **[Ver implementação com Redis](https://github.com/JoaoGabriel-Lima/desafio-devops/tree/feat/redis?tab=readme-ov-file)**

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

### 3. Cache Redis (Implementação Alternativa)

**Implementação**: Cache distribuído usando Redis como armazenamento externo

**Configuração**:

```python
import redis

redis_client = Redis(host=os.environ.get('REDIS_HOST', 'localhost'), port=6379)
DEFAULT_TTL = 60  # Tempo de expiração padrão para cache em segundos

def static_text():
    cache_key = "texto_estatico"
    
    cached_text = redis_client.get(cache_key)
    if cached_text:
        print(f"CACHE HIT (Python): texto estático encontrado no cache Redis, tempo restante: {redis_client.ttl(cache_key)}")
        return Response(cached_text.decode('utf-8'), mimetype='text/plain')
    
    print("CACHE MISS (Python): texto estático não encontrado no cache Redis")
    redis_client.set(cache_key, "Texto estático (Python)", ex=DEFAULT_TTL)
    return Response("Texto estático (Python)", mimetype='text/plain')
```

**Características**:

- TTL: 1 minuto (60 segundos)
- Cache distribuído e persistente
- Suporte a expiraçãoautomática com comando `EXPIRE`
- Monitoramento de TTL com comando `TTL`
- Logs detalhados de CACHE HIT/MISS
- Configuração via variável de ambiente `REDIS_HOST`

## 📊 Observabilidade

### Métricas Coletadas

Ambas as aplicações expõem métricas no formato Prometheus:

**Contador de Requisições HTTP**:

```prometheus
http_total_requests{path="/", method="GET"}
```

### Dashboard Grafana

<img src="./assets/dash_grafana.png" alt="Dashboard Grafana" width="800"/>

- **Go App**: Taxa de requisições por minuto
- **Python App**: Taxa de requisições por minuto
- Configuração automática via provisioning

## 🗂️ Estrutura do Projeto

```text
desafio_devops/
├── .github/
│   └── workflows/
│       └── ci.yml                  # Pipeline CI/CD para Go e Python
├── .gitignore                      # Arquivos ignorados pelo Git
├── docker-compose.yml              # Orquestração da infraestrutura
├── README.md                       # Este arquivo
├── assets/
│   ├── dash_grafana.png           # Screenshot do dashboard Grafana
│   └── infra_arq.png              # Diagrama de arquitetura
├── app1_golang/                    # Aplicação Go (Standard Go Project Layout)
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/app/main.go
│   └── internal/
│       ├── cache/
│       │   ├── cache.go           # Implementação do cache em memória
│       │   └── cache_test.go      # Testes unitários do cache
│       └── server/
│           └── server.go          # Servidor HTTP com métricas
├── app2_python/                    # Aplicação Python
│   ├── Dockerfile
│   ├── requirements.txt
│   ├── app.py                     # Flask app com métricas
│   └── test_app.py                # Testes unitários Python
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

Este projeto implementa um fluxo de atualização de código que automatiza o processo de testes para os microsserviços. Uma pipeline foi projetada para lidar com cenários onde desenvolvedores fazem alterações em qualquer uma das aplicações (Go ou Python) e acionam automaticamente a pipeline de CI.

### 🚀 Processo Detalhado de Atualização

### 🚀 Processo de Atualização

#### 1. Desenvolvimento Local

- Desenvolvedor identifica um bug ou implementa nova feature
- Realiza alterações no código da aplicação Go (`app1_golang/`) ou Python (`app2_python/`)
- Executa testes localmente para validação inicial

#### 2. Versionamento e Trigger da Pipeline

- Commit das alterações com mensagem descritiva
- Push para branch `main` ou criação de Pull Request
- Pipeline CI é automaticamente acionada via GitHub Actions

#### 3. Detecção Inteligente de Mudanças

- Utiliza `dorny/paths-filter@v3` para detectar quais aplicações foram modificadas
- Executa jobs específicos apenas para as aplicações alteradas:
  - Se mudanças em `app1_golang/**` → executa job `teste-e-build-go`
  - Se mudanças em `app2_python/**` → executa job `teste-e-build-python`

#### 4. Validação Automatizada

- **Testes Unitários**: Executados automaticamente para garantir qualidade
  - Go: `go test -v ./...`
  - Python: `pytest`
- **Análise de Segurança**: Trivy identifica vulnerabilidades críticas
- **Build**: Compilação/validação da aplicação

#### 5. Falha e Recuperação

- Pipeline falha se testes não passarem ou vulnerabilidades críticas forem encontradas
- Desenvolvedor recebe notificação automática
- Correções podem ser feitas e pipeline re-executada

### 📊 Diagrama do Fluxo de Atualização
<img src="./assets/fluxo_atualizacao.png" alt="Fluxo de Atualização de Componentes" width="800"/>

### 🛠️ Configuração da Pipeline

A pipeline está configurada no arquivo `.github/workflows/ci.yml` e inclui:

- **Triggers**: Push para `main` e Pull Requests
- **Path Filtering**: Execução condicional baseada nos arquivos alterados
- **Paralelização**: Jobs independentes para Go e Python
- **Security First**: Análise de vulnerabilidades com exit code 1 para críticas
- **Feedback Rápido**: Notificações imediatas em caso de falha


## 🚀 Pontos de Melhoria

### Melhorias do Aplicativo Go 
### Melhorias do Aplicativo Python
### Melhorias Gerais da Infraestrutura


## 🧪 Testes

### Testando Cache do Go (10 segundos)

```bash
# Primeira requisição (CACHE MISS)
curl http://localhost/go/time

# Segunda requisição imediata (CACHE HIT)
curl http://localhost/go/time

# Aguardar 10 segundos e fazer nova requisição (CACHE MISS)
sleep 11
curl http://localhost/go/time
```

### Testando Cache do Nginx (1 minuto)

```bash
# Primeira requisição (CACHE MISS)
curl -H "X-Cache-Status: debug" http://localhost/python/time

# Segunda requisição imediata (CACHE HIT)
curl -H "X-Cache-Status: debug" http://localhost/python/time

# Verificar status do cache no header
curl -I http://localhost/python/time
```

### Validando Métricas

```bash
# Métricas da aplicação Go
curl http://localhost/go/metrics

# Métricas da aplicação Python
curl http://localhost/python/metrics

# Verificar Prometheus
open http://localhost:9090

# Verificar Grafana
open http://localhost:3000
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
