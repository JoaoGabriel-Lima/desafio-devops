import os
import time
from flask import Flask
from redis import Redis

app = Flask(__name__)
redis_client = Redis(host=os.environ.get('REDIS_HOST', 'localhost'), port=6379)

DEFAULT_TTL = 60  # Tempo de expiração padrão para cache em segundos


@app.route("/")
def index():
    return "Texto estático (Python)"


@app.route("/time")
def get_time():
    cache_key = "hora_atual"

    cached_time = redis_client.get(cache_key)
    if cached_time:
        print("CACHE HIT (Python): hora atual encontrada no cache Redis, tempo restante para expiração:",
              redis_client.ttl(cache_key))
        return f"Hora atual guardada na cache do Redis: {cached_time.decode('utf-8')}"

    print("CACHE MISS (Python): hora atual não encontrada no cache Redis, gerando nova hora")
    current_time = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())

    redis_client.set(cache_key, current_time, ex=DEFAULT_TTL)

    return f"Hora atual do servidor (Python): {current_time}"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
