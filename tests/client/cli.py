import requests
import json
import os
from getpass import getpass
from datetime import datetime
import sys
from rich.console import Console
from rich.table import Table
from rich.prompt import Prompt, Confirm

class BankCLI:
    def __init__(self):
        self.base_url = "http://localhost:8080"
        self.token = None
        self.console = Console()
        self.config_file = "bank_config.json"
        self.user_roles = []
        self.load_config()

    def load_config(self):
        if os.path.exists(self.config_file):
            with open(self.config_file, 'r') as f:
                config = json.load(f)
                self.token = config.get('token')
                self.user_roles = config.get('roles', [])

    def save_config(self):
        with open(self.config_file, 'w') as f:
            json.dump({
                'token': self.token,
                'roles': self.user_roles
            }, f)

    def is_admin(self):
        return 'admin' in self.user_roles

    def clear_screen(self):
        os.system('clear' if os.name == 'posix' else 'cls')

    def print_header(self):
        self.clear_screen()
        self.console.print("[bold blue]Банковское приложение[/bold blue]")
        self.console.print("=" * 50)

    def register(self):
        self.print_header()
        self.console.print("[bold]Регистрация нового пользователя[/bold]")
        
        fio = Prompt.ask("Введите ФИО")
        username = Prompt.ask("Введите имя пользователя")
        email = Prompt.ask("Введите email")
        password = getpass("Введите пароль")
        
        try:
            response = requests.post(
                f"{self.base_url}/api/auth/register",
                json={
                    "fio": fio,
                    "username": username,
                    "email": email,
                    "password": password
                }
            )
            
            if response.status_code == 201:
                self.console.print("[green]Регистрация успешна![/green]")
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def login(self):
        self.print_header()
        self.console.print("[bold]Вход в систему[/bold]")
        
        username = Prompt.ask("Введите имя пользователя")
        password = getpass("Введите пароль")
        
        try:
            response = requests.post(
                f"{self.base_url}/api/auth/login",
                json={
                    "username": username,
                    "password": password
                }
            )
            
            if response.status_code == 200:
                data = response.json()
                self.token = data.get('token')
                self.user_roles = data.get('roles', [])
                self.save_config()
                self.console.print("[green]Вход выполнен успешно![/green]")
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def create_account(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header()
        self.console.print("[bold]Создание нового счета[/bold]")
        
        name = Prompt.ask("Введите название счета")
        currency = Prompt.ask("Введите валюту", default="RUB")
        
        try:
            response = requests.post(
                f"{self.base_url}/api/accounts",
                headers={"Authorization": f"Bearer {self.token}"},
                json={
                    "name": name,
                    "currency": currency
                }
            )
            
            if response.status_code == 201:
                self.console.print("[green]Счет создан успешно![/green]")
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def create_card(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header()
        self.console.print("[bold]Создание новой карты[/bold]")
        
        # Получаем список счетов
        try:
            response = requests.get(
                f"{self.base_url}/api/accounts",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if response.status_code != 200:
                self.console.print(f"[red]Ошибка при получении счетов: {response.json().get('message')}[/red]")
                return
            
            accounts = response.json().get('accounts', [])
            if not accounts:
                self.console.print("[red]У вас нет доступных счетов[/red]")
                return
            
            # Показываем таблицу счетов
            table = Table(show_header=True, header_style="bold magenta")
            table.add_column("ID")
            table.add_column("Название")
            table.add_column("Баланс")
            
            for account in accounts:
                table.add_row(
                    str(account['id']),
                    account['name'],
                    str(account['balance'])
                )
            
            self.console.print(table)
            
            account_id = Prompt.ask("Выберите ID счета")
            
            # Создаем карту
            response = requests.post(
                f"{self.base_url}/api/cards",
                headers={"Authorization": f"Bearer {self.token}"},
                json={
                    "account_id": int(account_id)
                }
            )
            
            if response.status_code == 201:
                self.console.print("[green]Карта создана успешно![/green]")
                card_data = response.json().get('card', {})
                
                # Показываем данные карты
                card_table = Table(show_header=True, header_style="bold green")
                card_table.add_column("Поле")
                card_table.add_column("Значение")
                
                for key, value in card_data.items():
                    card_table.add_row(key, str(value))
                
                self.console.print(card_table)
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def manage_users(self):
        if not self.is_admin():
            self.console.print("[red]Ошибка: Недостаточно прав[/red]")
            return
        
        self.print_header()
        self.console.print("[bold]Управление пользователями[/bold]")
        
        try:
            response = requests.get(
                f"{self.base_url}/api/admin/users",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if response.status_code == 200:
                users = response.json().get('users', [])
                
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("ID")
                table.add_column("ФИО")
                table.add_column("Имя пользователя")
                table.add_column("Email")
                table.add_column("Роли")
                
                for user in users:
                    roles = ", ".join(user.get('roles', []))
                    table.add_row(
                        str(user['id']),
                        user['fio'],
                        user['username'],
                        user['email'],
                        roles
                    )
                
                self.console.print(table)
                
                action = Prompt.ask(
                    "Выберите действие",
                    choices=["assign_role", "remove_role", "back"],
                    default="back"
                )
                
                if action == "back":
                    return
                
                user_id = Prompt.ask("Введите ID пользователя")
                role = Prompt.ask("Введите роль")
                
                if action == "assign_role":
                    response = requests.post(
                        f"{self.base_url}/api/admin/users/{user_id}/roles",
                        headers={"Authorization": f"Bearer {self.token}"},
                        json={"role": role}
                    )
                else:
                    response = requests.delete(
                        f"{self.base_url}/api/admin/users/{user_id}/roles/{role}",
                        headers={"Authorization": f"Bearer {self.token}"}
                    )
                
                if response.status_code in [200, 201]:
                    self.console.print("[green]Операция выполнена успешно![/green]")
                else:
                    self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def manage_roles(self):
        if not self.is_admin():
            self.console.print("[red]Ошибка: Недостаточно прав[/red]")
            return
        
        self.print_header()
        self.console.print("[bold]Управление ролями[/bold]")
        
        try:
            response = requests.get(
                f"{self.base_url}/api/admin/roles",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if response.status_code == 200:
                roles = response.json().get('roles', [])
                
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("Название")
                table.add_column("Описание")
                
                for role in roles:
                    table.add_row(
                        role['name'],
                        role['description']
                    )
                
                self.console.print(table)
                
                action = Prompt.ask(
                    "Выберите действие",
                    choices=["create", "delete", "back"],
                    default="back"
                )
                
                if action == "back":
                    return
                
                if action == "create":
                    name = Prompt.ask("Введите название роли")
                    description = Prompt.ask("Введите описание роли")
                    
                    response = requests.post(
                        f"{self.base_url}/api/admin/roles",
                        headers={"Authorization": f"Bearer {self.token}"},
                        json={
                            "name": name,
                            "description": description
                        }
                    )
                else:
                    role_name = Prompt.ask("Введите название роли для удаления")
                    response = requests.delete(
                        f"{self.base_url}/api/admin/roles/{role_name}",
                        headers={"Authorization": f"Bearer {self.token}"}
                    )
                
                if response.status_code in [200, 201]:
                    self.console.print("[green]Операция выполнена успешно![/green]")
                else:
                    self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def get_all_cards(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header()
        self.console.print("[bold]Список ваших карт[/bold]")
        
        try:
            # Сначала получаем список счетов
            accounts_response = requests.get(
                f"{self.base_url}/api/accounts",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if accounts_response.status_code != 200:
                self.console.print(f"[red]Ошибка при получении счетов: {accounts_response.json().get('message')}[/red]")
                return
            
            accounts = accounts_response.json().get('accounts', [])
            accounts_dict = {acc['id']: {'name': acc['name'], 'balance': acc['balance']} for acc in accounts}
            
            # Затем получаем список карт
            cards_response = requests.get(
                f"{self.base_url}/api/cards",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if cards_response.status_code == 200:
                cards = cards_response.json().get('cards', [])
                if not cards:
                    self.console.print("[yellow]У вас пока нет карт[/yellow]")
                    return
                
                # Создаем таблицу для отображения карт
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("ID")
                table.add_column("Номер карты")
                table.add_column("Срок действия")
                table.add_column("Счет")
                table.add_column("Баланс")
                table.add_column("Дата создания")
                
                for card in cards:
                    account_info = accounts_dict.get(card['account_id'], {"name": "Неизвестный счет", "balance": 0})
                    
                    table.add_row(
                        str(card['id']),
                        card.get('number', '**** **** **** ****'),
                        card.get('expiry_date', '**/**'),
                        account_info['name'],
                        f"{account_info['balance']:.2f} ₽",
                        card.get('created_at', '')
                    )
                
                self.console.print(table)
            else:
                self.console.print(f"[red]Ошибка: {cards_response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def show_menu(self):
        while True:
            self.print_header()
            
            if self.token:
                role_text = ", ".join(self.user_roles) if self.user_roles else "нет ролей"
                self.console.print(f"[green]Вы авторизованы[/green] (Роли: {role_text})")
            else:
                self.console.print("[yellow]Вы не авторизованы[/yellow]")
            
            self.console.print("\n[bold]Меню:[/bold]")
            self.console.print("1. Регистрация")
            self.console.print("2. Вход")
            self.console.print("3. Создать счет")
            self.console.print("4. Создать карту")
            self.console.print("5. Мои карты")
            
            if self.is_admin():
                self.console.print("6. Управление пользователями")
                self.console.print("7. Управление ролями")
            
            self.console.print("0. Выход")
            
            choices = ["0", "1", "2", "3", "4", "5"]
            if self.is_admin():
                choices.extend(["6", "7"])
            
            choice = Prompt.ask("Выберите действие", choices=choices)
            
            if choice == "0":
                if Confirm.ask("Вы уверены, что хотите выйти?"):
                    sys.exit()
            elif choice == "1":
                self.register()
            elif choice == "2":
                self.login()
            elif choice == "3":
                self.create_account()
            elif choice == "4":
                self.create_card()
            elif choice == "5":
                self.get_all_cards()
            elif choice == "6" and self.is_admin():
                self.manage_users()
            elif choice == "7" and self.is_admin():
                self.manage_roles()

if __name__ == "__main__":
    cli = BankCLI()
    cli.show_menu() 