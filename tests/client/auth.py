import requests
import json
from rich.console import Console
from rich.prompt import Prompt
import time

class AuthManager:
    def __init__(self, base_url, console):
        self.base_url = base_url
        self.console = console
        self.token = None
        self.user_roles = []
        self.config_file = "bank_config.json"
        self.load_config()


    def load_config(self):
        try:
            with open(self.config_file, 'r') as f:
                config = json.load(f)
                self.token = config.get('token')
                self.user_roles = config.get('roles', [])
        except FileNotFoundError:
            pass


    def save_config(self):
        with open(self.config_file, 'w') as f:
            json.dump({
                'token': self.token,
                'roles': self.user_roles
            }, f)


    def register(self):
        self.console.print("\n[bold blue]Регистрация[/bold blue]")
        self.console.print("=" * 50)
        
        username = Prompt.ask("Введите имя пользователя")
        email = Prompt.ask("Введите email")
        fio = Prompt.ask("Введите ФИО (Иванов Иван Иванович или Ivanov Ivan Ivanovich)")
        password = Prompt.ask("Введите пароль", password=True)

        # Отладочная информация
        self.console.print(f"[yellow]Отправляем запрос на {self.base_url}/api/auth/register[/yellow]")
        self.console.print(f"[yellow]Данные: username={username}, email={email}, fio={fio}[/yellow]")

        try:
            response = requests.post(
                f"{self.base_url}/api/auth/register",
                json={
                    "username": username,
                    "email": email,
                    "fio": fio,
                    "password": password
                }
            )

            # Отладочная информация
            self.console.print(f"[yellow]Статус ответа: {response.status_code}[/yellow]")
            self.console.print(f"[yellow]Ответ сервера: {response.text}[/yellow]")

            if response.status_code == 429:
                self.console.print("[yellow]Слишком много запросов. Подождите немного...[/yellow]")
                time.sleep(2)  # Ждем 2 секунды
                return self.register()  # Пробуем снова

            if response.status_code == 201:
                self.console.print("[green]Регистрация успешна![/green]")
            else:
                error_message = response.json().get('message')
                if error_message is None:
                    error_message = response.text
                self.console.print(f"[red]Ошибка: {error_message}[/red]")
        except Exception as e:
            self.console.print(f"[red]Ошибка при отправке запроса: {str(e)}[/red]")

    def login(self):
        self.console.print("\n[bold blue]Вход[/bold blue]")
        self.console.print("=" * 50)
        
        username = Prompt.ask("Введите имя пользователя")
        password = Prompt.ask("Введите пароль", password=True)

        response = requests.post(
            f"{self.base_url}/api/auth/login",
            json={
                "username": username,
                "password": password
            }
        )

        # Отладочная информация
        self.console.print(f"[yellow]Статус ответа: {response.status_code}[/yellow]")
        self.console.print(f"[yellow]Ответ сервера: {response.text}[/yellow]")

        if response.status_code == 200:
            data = response.json()
            self.token = data.get("token")
            self.user_roles = data.get('roles', [])
            self.save_config()
            self.console.print("[green]Вход выполнен успешно![/green]")
            # Отладочная информация о токене
            self.console.print(f"[yellow]Полученный токен: {self.token}[/yellow]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")


    def check_auth(self):
        if not self.token:
            self.console.print("[yellow]Токен отсутствует[/yellow]")
            return False
        
        try:
            response = requests.get(
                f"{self.base_url}/api/auth/auth-status",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            # Отладочная информация
            self.console.print(f"[yellow]Статус проверки авторизации: {response.status_code}[/yellow]")
            self.console.print(f"[yellow]Ответ сервера: {response.text}[/yellow]")
            
            if response.status_code == 429:
                self.console.print("[yellow]Слишком много запросов. Подождите немного...[/yellow]")
                time.sleep(2)  # Ждем 2 секунды
                return self.check_auth()  # Пробуем снова
            
            return response.status_code == 200
        except Exception as e:
            self.console.print(f"[red]Ошибка при проверке авторизации: {str(e)}[/red]")
            return False


    def logout(self):
        if not self.token:
            self.console.print("[yellow]Вы не авторизованы[/yellow]")
            return
        
        try:
            # Пытаемся отправить запрос на сервер для выхода
            response = requests.post(
                f"{self.base_url}/api/auth/logout",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            # Очищаем локальные данные независимо от ответа сервера
            self.token = None
            self.user_roles = []
            self.save_config()
            
            if response.status_code == 200:
                self.console.print("[green]Вы успешно вышли из системы[/green]")
            else:
                self.console.print("[yellow]Выход выполнен локально, но сервер не ответил корректно[/yellow]")
        except Exception as e:
            self.console.print(f"[yellow]Ошибка при выходе: {str(e)}[/yellow]")
            # Всё равно очищаем локальные данные
            self.token = None
            self.user_roles = []
            self.save_config()


    def is_admin(self):
        return 'admin' in self.user_roles 