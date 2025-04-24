import unittest
import http.client
import json
import logging
import sys

from utils import Storage

LOGIN_USER = "12312312332233"
PASSWORD_USER = "testpassword"

# Настройка логирования
logger = logging.getLogger("test")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")


class BaseAPITest(unittest.TestCase):
    def setUp(self):
        self.storage = Storage() 
        self.host = "localhost:8080"
        self.conn = http.client.HTTPConnection(self.host)
        self.login_user = LOGIN_USER
        self.login_password = PASSWORD_USER
        self.jwt_token = ""
        self.headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': f'Bearer {self.storage.get("token")}'
        }

    def send_request(self, method, path, body=None, headers=None):
        self.conn.request(method, path, body, headers or {})
        response = self.conn.getresponse()
        data = response.read()
        return response.status, data

    def tearDown(self):
        # Закрыть соединение после тестов
        self.conn.close()



class TestAuthAPIRequests(BaseAPITest):
    def setUp(self):
        super().setUp()

    def send_request(self, method, path, body=None, headers=None):
        self.conn.request(method, path, body, headers or {})
        response = self.conn.getresponse()
        data = response.read()
        return response.status, data

    def test_post_register_request(self):
        payload = json.dumps({
            "username": self.login_user,
            "password": self.login_password
        })
        try:
            status, data = self.send_request("POST", "/api/auth/register", body=payload, headers=self.headers)
            response_data = json.loads(data)

            self.assertEqual(status, 201, f"Ожидаемый статус: 201, полученный: {status}")
            self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")

            logger.info("Тест регистрации завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест регистрации не прошел: {e}")
            raise
    
    def test_post_login_request(self):
        payload = json.dumps({
            "username": self.login_user,
            "password": self.login_password,
        })
        try:
            status, data = self.send_request("POST", "/api/auth/login", body=payload)
            response_data = json.loads(data)

            self.assertEqual(status, 200, f"Ожидаемый статус: 200, полученный: {status}")
            self.token = response_data.get("token", "")
            self.headers['Authorization'] = f'Bearer {self.token}'
            self.storage.set("token", self.token)
            self.assertTrue(self.token, "Токен авторизации отсутствует в ответе")

            logger.info("Тест авторизации завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест авторизации не прошел: {e}")
            raise

    def tearDown(self):
        # Закрыть соединение после тестов
        self.conn.close()

class TestAccountAPIRequests(BaseAPITest):
    def setUp(self):
        super().setUp()

    def send_request(self, method, path, body=None, headers=None):
        self.conn.request(method, path, body, headers or {})
        response = self.conn.getresponse()
        data = response.read()
        return response.status, data


    def test_post_create_account_request(self):
        payload = json.dumps({
            "name": "Test Account",
        })
        try:
            status, data = self.send_request("POST", "/api/accounts", body=payload, headers=self.headers)
            response_data = json.loads(data)

            self.assertEqual(status, 201, f"Ожидаемый статус: 201, полученный: {status}")
            self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")

            logger.info("Тест создания счета завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест создания счета не прошел: {e}")
            raise

    def test_get_accounts_request(self):
        try:
            status, data = self.send_request("GET", "/api/accounts", headers=self.headers)
            response_data = json.loads(data)

            self.assertEqual(status, 200, f"Ожидаемый статус: 200, полученный: {status}")
            self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")

            logger.info("Тест получения списка счетов завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест получения списка счетов не прошел: {e}")
            raise


    


if __name__ == "__main__":
    suite = unittest.TestSuite()
    suite.addTest(TestAuthAPIRequests("test_post_register_request"))
    suite.addTest(TestAuthAPIRequests("test_post_login_request"))
    suite.addTest(TestAccountAPIRequests("test_post_create_account_request"))
    suite.addTest(TestAccountAPIRequests("test_get_accounts_request"))

    # runner = QuietTestRunner()
    # runner.run(suite)
    runner = unittest.TextTestRunner()
    runner.run(suite)