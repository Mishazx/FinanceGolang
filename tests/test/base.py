import logging
from unittest import TestCase
from utils import Storage
import http.client

# Настройка логирования
logger = logging.getLogger("test")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")


class BaseAPITest(TestCase):
    def setUp(self):
        self.storage = Storage() 
        self.host = "localhost:8080"
        self.conn = http.client.HTTPConnection(self.host)
        self.headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': f'Bearer {self.storage.get("token")}'
        }
        self.users = [
            {
                "fio": "Ivanov Ivan Ivanovich", 
                "username": "user1", 
                "password": "password1", 
                "email": "6mNlD@example.com"
            },
            {
                "fio": "Petrov Petr Petrovich",
                "username": "user2",
                "password": "password2", 
                "email": "B0P0o@example.com"
            },
        ]
        self.tokens = {}

    def send_request(self, method, path, body=None, headers=None):
        self.conn.request(method, path, body, headers or {})
        response = self.conn.getresponse()
        data = response.read()
        return response.status, data

    def tearDown(self):
        # Закрыть соединение после тестов
        self.conn.close()
