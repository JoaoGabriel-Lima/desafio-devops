package cache

import (
	"testing"
	"time"
)

const (
	msgItemDeveriaTerSidoEncontrado = "Item deveria ter sido encontrado no cache"
	msgValorEsperadoObtido = "Valor esperado: %s, valor obtido: %s"
)

func TestCacheSetEGet(t *testing.T) {
	cache := NewCache()
	chave := "teste_chave"
	valor := "teste_valor"
	ttl := 5 * time.Second

	// Testa o método Set
	cache.Set(chave, valor, ttl)

	// Verifica se o item foi adicionado ao cache usando Get
	valorRetornado, expiracao, encontrado := cache.Get(chave)

	if !encontrado {
		t.Error(msgItemDeveriaTerSidoEncontrado)
	}

	if valorRetornado != valor {
		t.Errorf(msgValorEsperadoObtido, valor, valorRetornado)
	}

	if expiracao <= time.Now().UnixNano() {
		t.Error("Expiração deveria ser no futuro")
	}
}

func TestCacheGetMultiplosItens(t *testing.T) {
	cache := NewCache()
	
	// Adiciona múltiplos itens
	itens := map[string]string{
		"chave1": "valor1",
		"chave2": "valor2",
		"chave3": "valor3",
	}
	
	ttl := 5 * time.Second
	
	for chave, valor := range itens {
		cache.Set(chave, valor, ttl)
	}

	// Verifica se todos os itens podem ser recuperados
	for chave, valorEsperado := range itens {
		valorRetornado, _, encontrado := cache.Get(chave)

		if !encontrado {
			t.Errorf("Item com chave %s deveria ter sido encontrado no cache", chave)
		}

		if valorRetornado != valorEsperado {
			t.Errorf(msgValorEsperadoObtido, valorEsperado, valorRetornado)
		}
	}
}

func TestCacheGetChaveInexistente(t *testing.T) {
	cache := NewCache()
	chave := "chave_inexistente"

	// Tenta buscar uma chave que não existe
	_, _, encontrado := cache.Get(chave)

	if encontrado {
		t.Error("Item não deveria ter sido encontrado no cache")
	}
}

func TestCacheExpiracao(t *testing.T) {
	cache := NewCache()
	chave := "teste_chave"
	valor := "teste_valor"
	ttl := 100 * time.Millisecond // TTL muito baixo para teste

	// Adiciona item ao cache
	cache.Set(chave, valor, ttl)

	// Verifica se o item existe inicialmente
	valorRetornado, _, encontrado := cache.Get(chave)
	if !encontrado {
		t.Error(msgItemDeveriaTerSidoEncontrado)
	}
	if valorRetornado != valor {
		t.Errorf(msgValorEsperadoObtido, valor, valorRetornado)
	}

	// Espera o TTL expirar
	time.Sleep(150 * time.Millisecond)

	// Verifica se o item foi removido após expiração
	_, _, encontradoAposExpiracao := cache.Get(chave)
	if encontradoAposExpiracao {
		t.Error("Item não deveria ser encontrado após expiração")
	}
}

func TestCacheConcorrencia(t *testing.T) {
	cache := NewCache()
	chave := "teste_concorrencia"
	valor := "valor_concorrencia"
	ttl := 5 * time.Second

	// Testa operações concorrentes
	done := make(chan bool, 2)

	// Goroutine 1: Escreve no cache
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(chave, valor, ttl)
		}
		done <- true
	}()

	// Goroutine 2: Lê do cache
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get(chave)
		}
		done <- true
	}()

	// Espera ambas as goroutines terminarem
	<-done
	<-done

	// Verifica se o cache ainda funciona após operações concorrentes
	valorFinal, _, encontrado := cache.Get(chave)
	if !encontrado {
		t.Error("Item deveria estar no cache após operações concorrentes")
	}
	if valorFinal != valor {
		t.Errorf(msgValorEsperadoObtido, valor, valorFinal)
	}
}

func TestCacheInterface(t *testing.T) {
	// Testa se o CacheEmMemoria implementa a interface corretamente
	var cache Interface = NewCache()
	
	chave := "teste_interface"
	valor := "valor_interface"
	ttl := 5 * time.Second

	cache.Set(chave, valor, ttl)
	valorRetornado, _, encontrado := cache.Get(chave)

	if !encontrado {
		t.Error(msgItemDeveriaTerSidoEncontrado)
	}

	if valorRetornado != valor {
		t.Errorf(msgValorEsperadoObtido, valor, valorRetornado)
	}
}
