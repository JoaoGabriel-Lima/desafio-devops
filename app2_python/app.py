import time
from flask import Flask, Response

app = Flask(__name__)


@app.route("/")
def index():
    return Response("Desafio DevOps - João Gabriel Lima Marinho - Servidor Python", mimetype='text/plain')


@app.route("/static-text")
def static_text():
    return Response("Texto estático (Python)", mimetype='text/plain')


@app.route("/time")
def get_time():
    current_time = time.strftime('%d/%m/%Y | %H:%M:%S', time.gmtime())
    return Response(f"Hora atual do servidor (Python): {current_time}", mimetype='text/plain')


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
