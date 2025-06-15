import time
from flask import Flask, Response, request
from prometheus_client import Counter, make_wsgi_app
from werkzeug.middleware.dispatcher import DispatcherMiddleware


app = Flask(__name__)

TOTAL_REQUESTS = Counter(
    "http_total_requests",
    "Total de requisições HTTP recebidas",
    ["path", "method"]
)


@app.route("/")
def index():
    TOTAL_REQUESTS.labels(path="/", method=request.method).inc()
    return Response("Desafio DevOps - João Gabriel Lima Marinho - Servidor Python", mimetype='text/plain')


@app.route("/static-text")
def static_text():
    TOTAL_REQUESTS.labels(path="/static-text", method=request.method).inc()
    return Response("Texto estático (Python)", mimetype='text/plain')


@app.route("/time")
def get_time():
    TOTAL_REQUESTS.labels(path="/time", method=request.method).inc()
    current_time = time.strftime('%d/%m/%Y | %H:%M:%S', time.gmtime())
    return Response(f"Hora atual do servidor (Python): {current_time}", mimetype="text/plain")


app.wsgi_app = DispatcherMiddleware(app.wsgi_app, {
    '/metrics': make_wsgi_app()
})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
