import unittest
import json
import logging
import sys

from test_account import TestAccountAPIRequests
from test_auth import TestAuthAPIRequests
# from myrandom import generation_name_account

from base import BaseAPITest, logger

class TestCardAPIRequests(BaseAPITest):
    def setUp(self):
        super().setUp()

    def test_post_card_request(self):
        payload = json.dumps({
            "name": "Test Card",
        })
        try:
            status, data = self.send_request("POST", "/api/cards", body=payload, headers=self.headers)
            response_data = json.loads(data)

            self.assertEqual(status, 201, f"Ожидаемый статус: 201, полученный: {status}")
            self.assertEqual(response_data["status"], "success", f"Ожидаемый статус ответа: 'success', полученный: {response_data.get('status')}")
            logger.info("Тест создания карты завершен успешно")

        except AssertionError as e:
            logger.error(f"Тест создания карты не прошел: {e}")
            raise
    


if __name__ == "__main__":
    suite = unittest.TestSuite()
    # Добавление тестов для пользователя smirnov
    suite.addTest(TestAuthAPIRequests("test_register_user1"))
    suite.addTest(TestAuthAPIRequests("test_login_user1"))
    # Добавление тестов для пользователя ivanov
    suite.addTest(TestAuthAPIRequests("test_register_user2"))
    suite.addTest(TestAuthAPIRequests("test_login_user2"))


    # регистрация двух счетов и вывод списка счетов для user1
    suite.addTest(TestAccountAPIRequests("test_create_account_user1")) 
    suite.addTest(TestAccountAPIRequests("test_get_account_user1"))

    # регистрация счета и вывод списка счетов для user2
    suite.addTest(TestAccountAPIRequests("test_create_account_user2"))
    suite.addTest(TestAccountAPIRequests("test_get_account_user2"))
    # регистрация счета и вывод списка счетов для иванова

    # регистрация карты и вывод списка карт


    # runner = QuietTestRunner()
    # runner.run(suite)
    runner = unittest.TextTestRunner()
    runner.run(suite)