import json
import unittest
from base import BaseAPITest, logger

class TestAuthAPIRequests(BaseAPITest):
    def setUp(self):
        super().setUp()

    def test_register_user1(self):
        user = self.users[0]  # Access the first user
        self.register_user(user)

    def test_register_user2(self):
        user = self.users[1]  # Access the second user
        self.register_user(user)

    def register_user(self, user):
        payload = json.dumps({
            "username": user["username"],
            "password": user["password"],
            "email": user["email"]
        })
        try:
            status, data = self.send_request("POST", "/api/auth/register", body=payload, headers=self.headers)
            response_data = json.loads(data)

            self.assertEqual(status, 201, f"Ожидаемый статус: 201, полученный: {status}")
            self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")

            logger.info(f"Тест регистрации для пользователя {user['username']} завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест регистрации для пользователя {user['username']} не прошел: {e}")
            raise

    def test_login_user1(self):
        user = self.users[0]  # Access the first user
        self.login_user(user)

    def test_login_user2(self):
        user = self.users[1]  # Access the second user
        self.login_user(user)

    def login_user(self, user):
        payload = json.dumps({
            "username": user["username"],
            "password": user["password"],
        })
        try:
            status, data = self.send_request("POST", "/api/auth/login", body=payload)
            response_data = json.loads(data)

            self.assertEqual(status, 200, f"Ожидаемый статус: 200, полученный: {status}")
            token = response_data.get("token", "")
            self.headers['Authorization'] = f'Bearer {token}'
            self.storage.set(f"token_{user['username']}", token)
            self.assertTrue(token, "Токен авторизации отсутствует в ответе")

            logger.info(f"Тест авторизации для пользователя {user['username']} завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест авторизации для пользователя {user['username']} не прошел: {e}")
            raise

    def tearDown(self):
        # Закрыть соединение после тестов
        self.conn.close()
