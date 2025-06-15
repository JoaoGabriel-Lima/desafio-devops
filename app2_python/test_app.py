import unittest
import time
from app import app

# Testes unitários para o App 2


class TestApp(unittest.TestCase):
    def setUp(self):
        """Configuração inicial para cada teste"""
        self.app = app.test_client()
        self.app.testing = True

    def test_index_route(self):
        """Testa a rota principal /"""
        response = self.app.get('/')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Desafio DevOps', response.data)
        self.assertIn(b'Servidor Python', response.data)

    def test_static_text_route(self):
        """Testa a rota /static-text"""
        response = self.app.get('/static-text')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Texto est\xc3\xa1tico (Python)', response.data)

    def test_time_route(self):
        """Testa a rota /time"""
        response = self.app.get('/time')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Hora atual do servidor (Python)', response.data)
        # Verifica se contém um formato de data válido
        self.assertRegex(response.data.decode(),
                         r'\d{2}/\d{2}/\d{4} \| \d{2}:\d{2}:\d{2}')

    def test_metrics_endpoint(self):
        """Testa se o endpoint /metrics está disponível"""
        response = self.app.get('/metrics')
        self.assertEqual(response.status_code, 200)
        # Verifica se contém métricas do Prometheus
        self.assertIn(b'http_total_requests', response.data)

    def test_nonexistent_route(self):
        """Testa uma rota que não existe"""
        response = self.app.get('/rota-inexistente')
        self.assertEqual(response.status_code, 404)

    def test_content_type(self):
        """Testa se o content-type está correto"""
        response = self.app.get('/')
        self.assertEqual(response.content_type, 'text/plain; charset=utf-8')

    def test_multiple_requests_counter(self):
        """Testa se o contador de requisições está funcionando"""
        # Faz algumas requisições
        self.app.get('/')
        self.app.get('/static-text')
        self.app.get('/time')

        # Verifica se as métricas foram atualizadas
        response = self.app.get('/metrics')
        self.assertIn(b'http_total_requests', response.data)


if __name__ == '__main__':
    unittest.main()
