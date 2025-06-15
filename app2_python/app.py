import os
import time
import logging
from flask import Flask, Response
from redis import Redis

app = Flask(__name__)
redis_client = Redis(host=os.environ.get('REDIS_HOST', 'localhost'), port=6379)

DEFAULT_TTL = 60  # Tempo de expiração padrão para cache em segundos


@app.route("/")
def index():
    return Response("Desafio DevOps - João Gabriel Lima Marinho - Servidor Python", mimetype='text/plain')


@app.route("/static-text")
def static_text():
    cache_key = "texto_estatico"

    cached_text = redis_client.get(cache_key)
    if cached_text:
        print("CACHE HIT (Python): texto estático encontrado no cache Redis, tempo restante para expiração:",
              redis_client.ttl(cache_key))
        return Response(cached_text.decode('utf-8'), mimetype='text/plain')

    print("CACHE MISS (Python): texto estático não encontrado no cache Redis, gerando novo texto")
    redis_client.set(cache_key, "Texto estático (Python)", ex=DEFAULT_TTL)
    print("Texto estático armazenado no cache Redis")
    return Response("Texto estático (Python)", mimetype='text/plain')


@app.route("/time")
def get_time():
    cache_key = "hora_atual"

    cached_time = redis_client.get(cache_key)
    if cached_time:
        print("CACHE HIT (Python): hora atual encontrada no cache Redis, tempo restante para expiração:",
              redis_client.ttl(cache_key))
        return Response(f"Hora atual guardada na cache do Redis: {cached_time.decode('utf-8')}", mimetype='text/plain')

    print("CACHE MISS (Python): hora atual não encontrada no cache Redis, gerando nova hora")
    current_time = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())

    redis_client.set(cache_key, current_time, ex=DEFAULT_TTL)

    return Response(f"Hora atual do servidor (Python): {current_time}", mimetype='text/plain')


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
