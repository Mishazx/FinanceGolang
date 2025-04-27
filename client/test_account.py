import json
from base import BaseAPITest, logger


class TestAccountAPIRequests(BaseAPITest):
    def setUp(self):
        super().setUp()
        self.token_user1 = self.storage.get("token_user1")
        self.token_user2 = self.storage.get("token_user2")
        self.headers_user1 = {"Authorization": f"Bearer {self.token_user1}"}
        self.headers_user2 = {"Authorization": f"Bearer {self.token_user2}"}
        # logger.info("Начало тестирования API счетов")

    def create_account(self, user, account_name, headers):
        payload = json.dumps({"name": account_name})
        status, data = self.send_request("POST", "/api/accounts", body=payload, headers=headers)
        response_data = json.loads(data)

        self.assertEqual(status, 201, f"Ожидаемый статус: 201, полученный: {status}")
        self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")
        logger.info(f"Тест создания счета для {user} завершен успешно: {account_name}")

    def get_accounts(self, user, headers):
        status, data = self.send_request("GET", "/api/accounts", headers=headers)
        response_data = json.loads(data)

        self.assertEqual(status, 200, f"Ожидаемый статус: 200, полученный: {status}")
        self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")
        logger.info(f"Тест получения списка счетов для {user} завершен успешно")

    def test_create_account_user1(self):
        # Create two accounts for user1
        for i in range(2):
            self.create_account("user1", f"Test Account User1 {i + 1}", self.headers_user1)

    def test_create_account_user2(self):
        # Create one account for user2
        self.create_account("user2", "Test Account User2", self.headers_user2)

    def test_get_account_user1(self):
        # Get accounts for user1
        self.get_accounts("user1", self.headers_user1)

    def test_get_account_user2(self):
        # Get accounts for user2
        self.get_accounts("user2", self.headers_user2)

    def tearDown(self):
        # Закрыть соединение после тестов
        self.conn.close()
