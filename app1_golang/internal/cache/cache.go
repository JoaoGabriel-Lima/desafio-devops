package cache

import (
	"log"
	"sync"
	"time"
)

type Interface interface {
	Set(chave string, valor string, ttl time.Duration)
	Get(chave string) (string, int64, bool)
}

type cacheItem struct {
	valor      string
	expiracao  int64
}

type CacheEmMemoria struct {
	sync.RWMutex
	items map[string]cacheItem
}

func NewCache() *CacheEmMemoria {
	return &CacheEmMemoria{
		items: make(map[string]cacheItem),
	}
}

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