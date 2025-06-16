package cache

import (
	"log"
	"sync"
	"time"
)

// Define os métodos que a interface deve implementar
type Interface interface {
	Set(chave string, valor string, ttl time.Duration)
	Get(chave string) (string, int64, bool)
}

// cacheItem representa um item armazenado no cache, incluindo o valor e a data de expiração
type cacheItem struct {
	valor      string
	expiracao  int64
}

// CacheEmMemoria é uma implementação da interface Interface que armazena dados em memória
type CacheEmMemoria struct {
	sync.RWMutex
	items map[string]cacheItem
}

// NewCache cria uma nova instância de CacheEmMemoria
func NewCache() *CacheEmMemoria {
	return &CacheEmMemoria{
		items: make(map[string]cacheItem),
	}
}

// Set armazena um valor no cache com uma chave e um tempo de vida (TTL) especificados
func (cm *CacheEmMemoria) Set(chave string, valor string, ttl time.Duration) {
	cm.Lock()
	defer cm.Unlock()

	expiracao := time.Now().Add(ttl).UnixNano()
	cm.items[chave] = cacheItem{
		valor:     valor,
		expiracao: expiracao,
	}
	log.Printf("CACHE SET: chave=%s, valor=%s, ttl=%s", chave, valor, ttl)
}

// Get recupera um valor do cache usando uma chave, retornando o valor, a data de expiração e um booleano indicando se o item foi encontrado
func (cm *CacheEmMemoria) Get(chave string) (string, int64, bool) {
	cm.RLock()
	defer cm.RUnlock()

	var item, encontrado = cm.items[chave]
	if !encontrado {
		return "", 0, false
	}

	if time.Now().UnixNano() > item.expiracao {
		delete(cm.items, chave)
		return "", 0, false
	}

	return item.valor, item.expiracao, true
}