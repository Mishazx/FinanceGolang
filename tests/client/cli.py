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

    def print_header(self, title):
        self.clear_screen()
        self.console.print(f"\n[bold blue]{title}[/bold blue]")
        self.console.print("=" * 50)

    def register(self):
        self.print_header("Регистрация")
        username = Prompt.ask("Введите имя пользователя")
        email = Prompt.ask("Введите email")
        password = Prompt.ask("Введите пароль", password=True)

        response = requests.post(
            f"{self.base_url}/api/auth/register",
            json={
                "username": username,
                "email": email,
                "password": password
            }
        )

        if response.status_code == 201:
            self.console.print("[green]Регистрация успешна![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def login(self):
        self.print_header("Вход")
        username = Prompt.ask("Введите имя пользователя")
        password = Prompt.ask("Введите пароль", password=True)

        response = requests.post(
            f"{self.base_url}/api/auth/login",
            json={
                "username": username,
                "password": password
            }
        )

        if response.status_code == 200:
            data = response.json()
            self.token = data.get("token")
            self.user_roles = data.get('roles', [])
            self.save_config()
            self.console.print("[green]Вход выполнен успешно![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def create_account(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("Создание счета")
        name = Prompt.ask("Введите название счета")

        response = requests.post(
            f"{self.base_url}/api/accounts",
            headers={"Authorization": f"Bearer {self.token}"},
            json={"name": name}
        )

        if response.status_code == 201:
            self.console.print("[green]Счет создан успешно![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def create_card(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header("Создание новой карты")
        
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
        
        self.print_header("Управление пользователями")
        
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
        
        self.print_header("Управление ролями")
        
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
        
        self.print_header("Список ваших карт")
        
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

    def deposit(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("Пополнение счета")
        account_id = Prompt.ask("Введите ID счета")
        amount = float(Prompt.ask("Введите сумму"))
        description = Prompt.ask("Введите описание")

        response = requests.post(
            f"{self.base_url}/api/accounts/{account_id}/deposit",
            headers={"Authorization": f"Bearer {self.token}"},
            json={
                "amount": amount,
                "description": description
            }
        )

        if response.status_code == 200:
            self.console.print("[green]Счет успешно пополнен![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def withdraw(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("Снятие средств")
        account_id = Prompt.ask("Введите ID счета")
        amount = float(Prompt.ask("Введите сумму"))
        description = Prompt.ask("Введите описание")

        response = requests.post(
            f"{self.base_url}/api/accounts/{account_id}/withdraw",
            headers={"Authorization": f"Bearer {self.token}"},
            json={
                "amount": amount,
                "description": description
            }
        )

        if response.status_code == 200:
            self.console.print("[green]Средства успешно сняты![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def transfer(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("Перевод средств")
        from_account_id = Prompt.ask("Введите ID счета отправителя")
        to_account_id = Prompt.ask("Введите ID счета получателя")
        amount = float(Prompt.ask("Введите сумму"))
        description = Prompt.ask("Введите описание")

        response = requests.post(
            f"{self.base_url}/api/accounts/{from_account_id}/transfer",
            headers={"Authorization": f"Bearer {self.token}"},
            json={
                "to_account_id": int(to_account_id),
                "amount": amount,
                "description": description
            }
        )

        if response.status_code == 200:
            self.console.print("[green]Перевод выполнен успешно![/green]")
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def get_transactions(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("История транзакций")
        account_id = Prompt.ask("Введите ID счета")

        response = requests.get(
            f"{self.base_url}/api/accounts/{account_id}/transactions",
            headers={"Authorization": f"Bearer {self.token}"}
        )

        if response.status_code == 200:
            transactions = response.json().get("transactions", [])
            table = Table(show_header=True, header_style="bold magenta")
            table.add_column("Дата")
            table.add_column("Тип")
            table.add_column("Сумма")
            table.add_column("Описание")
            table.add_column("Статус")
            
            for transaction in transactions:
                table.add_row(
                    transaction["created_at"],
                    transaction["type"],
                    str(transaction["amount"]),
                    transaction["description"],
                    transaction["status"]
                )
            
            self.console.print(table)
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def check_auth(self):
        if not self.token:
            return False
        
        try:
            response = requests.get(
                f"{self.base_url}/api/auth/auth-status",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            return response.status_code == 200
        except:
            return False

    def logout(self):
        if not self.token:
            self.console.print("[yellow]Вы не авторизованы[/yellow]")
            return
        
        self.token = None
        self.user_roles = []
        self.save_config()
        self.console.print("[green]Вы успешно вышли из системы[/green]")
        input("\nНажмите Enter для продолжения...")

    def create_credit(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header("Оформление кредита")
        
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
            
            account_id = Prompt.ask("Выберите ID счета для получения кредита")
            amount = float(Prompt.ask("Введите сумму кредита"))
            term_months = int(Prompt.ask("Введите срок кредита в месяцах"))
            description = Prompt.ask("Введите описание кредита", default="Потребительский кредит")
            
            response = requests.post(
                f"{self.base_url}/api/credits",
                headers={"Authorization": f"Bearer {self.token}"},
                json={
                    "account_id": int(account_id),
                    "amount": amount,
                    "term_months": term_months,
                    "description": description
                }
            )
            
            if response.status_code == 201:
                credit_data = response.json().get('credit', {})
                self.console.print("[green]Кредит успешно оформлен![/green]")
                
                # Показываем информацию о кредите
                credit_table = Table(show_header=True, header_style="bold green")
                credit_table.add_column("Поле")
                credit_table.add_column("Значение")
                
                for key, value in credit_data.items():
                    credit_table.add_row(key, str(value))
                
                self.console.print(credit_table)
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def get_credits(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header("Мои кредиты")
        
        try:
            response = requests.get(
                f"{self.base_url}/api/credits",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if response.status_code == 200:
                credits = response.json().get('credits', [])
                if not credits:
                    self.console.print("[yellow]У вас пока нет кредитов[/yellow]")
                    return
                
                # Создаем таблицу для отображения кредитов
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("ID")
                table.add_column("Сумма")
                table.add_column("Ставка")
                table.add_column("Срок")
                table.add_column("Ежемесячный платеж")
                table.add_column("Статус")
                table.add_column("Дата начала")
                table.add_column("Дата окончания")
                
                for credit in credits:
                    table.add_row(
                        str(credit['id']),
                        f"{credit['amount']:.2f} ₽",
                        f"{credit['interest_rate']:.2f}%",
                        f"{credit['term_months']} мес.",
                        f"{credit['monthly_payment']:.2f} ₽",
                        credit['status'],
                        credit['start_date'],
                        credit['end_date']
                    )
                
                self.console.print(table)
                
                # Предлагаем посмотреть график платежей
                credit_id = Prompt.ask("Введите ID кредита для просмотра графика платежей (или Enter для выхода)")
                if credit_id:
                    self.get_payment_schedule(int(credit_id))
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def get_payment_schedule(self, credit_id):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header(f"График платежей по кредиту #{credit_id}")
        
        try:
            response = requests.get(
                f"{self.base_url}/api/credits/{credit_id}/schedule",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            
            if response.status_code == 200:
                schedule = response.json().get('schedule', [])
                if not schedule:
                    self.console.print("[yellow]Нет данных о платежах[/yellow]")
                    return
                
                # Создаем таблицу для отображения графика платежей
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("№")
                table.add_column("Дата")
                table.add_column("Основной долг")
                table.add_column("Проценты")
                table.add_column("Общая сумма")
                table.add_column("Статус")
                
                for payment in schedule:
                    table.add_row(
                        str(payment['payment_number']),
                        payment['payment_date'],
                        f"{payment['principal_amount']:.2f} ₽",
                        f"{payment['interest_amount']:.2f} ₽",
                        f"{payment['total_amount']:.2f} ₽",
                        payment['status']
                    )
                
                self.console.print(table)
                
                # Предлагаем оплатить платеж
                if Confirm.ask("Хотите оплатить платеж?"):
                    payment_number = int(Prompt.ask("Введите номер платежа"))
                    self.process_payment(credit_id, payment_number)
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def process_payment(self, credit_id, payment_number):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.print_header(f"Оплата платежа #{payment_number} по кредиту #{credit_id}")
        
        try:
            response = requests.post(
                f"{self.base_url}/api/credits/{credit_id}/payment",
                headers={"Authorization": f"Bearer {self.token}"},
                json={
                    "payment_number": payment_number
                }
            )
            
            if response.status_code == 200:
                self.console.print("[green]Платеж успешно обработан![/green]")
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
        
        input("\nНажмите Enter для продолжения...")

    def get_accounts(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.print_header("Мои счета")
        response = requests.get(
            f"{self.base_url}/api/accounts",
            headers={"Authorization": f"Bearer {self.token}"}
        )

        if response.status_code == 200:
            accounts = response.json().get("accounts", [])
            table = Table(show_header=True, header_style="bold magenta")
            table.add_column("ID")
            table.add_column("Название")
            table.add_column("Баланс")
            
            for account in accounts:
                table.add_row(
                    str(account["id"]),
                    account["name"],
                    f"{account['balance']:.2f} ₽"
                )
            
            self.console.print(table)
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")

    def show_menu(self):
        while True:
            self.print_header("Банковский терминал")
            
            is_authenticated = self.check_auth()
            
            if is_authenticated:
                role_text = ", ".join(self.user_roles) if self.user_roles else "нет ролей"
                self.console.print(f"[green]Вы авторизованы[/green] (Роли: {role_text})")
            else:
                self.console.print("[yellow]Вы не авторизованы[/yellow]")
                self.token = None
                self.user_roles = []
                self.save_config()
            
            self.console.print("\n[bold]Меню:[/bold]")
            
            if not is_authenticated:
                self.console.print("1. Регистрация")
                self.console.print("2. Вход")
                self.console.print("0. Завершить программу")
            else:
                self.console.print("1. Создать счет")
                self.console.print("2. Мои счета")
                self.console.print("3. Создать карту")
                self.console.print("4. Мои карты")
                self.console.print("5. Пополнить счет")
                self.console.print("6. Снять средства")
                self.console.print("7. Перевести средства")
                self.console.print("8. История транзакций")
                self.console.print("9. Оформить кредит")
                self.console.print("10. Мои кредиты")
                
                if self.is_admin():
                    self.console.print("a. Управление пользователями")
                    self.console.print("b. Управление ролями")
                
                self.console.print("l. Выйти из аккаунта")
                self.console.print("0. Завершить программу")
            
            if not is_authenticated:
                choices = ["0", "1", "2"]
            else:
                choices = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "l"]
                if self.is_admin():
                    choices.extend(["a", "b"])
            
            choice = Prompt.ask("Выберите действие", choices=choices)
            
            if choice == "0":
                if Confirm.ask("Вы уверены, что хотите завершить программу?"):
                    sys.exit()
            elif choice == "l" and is_authenticated:
                self.logout()
            elif choice == "1":
                if is_authenticated:
                    self.create_account()
                else:
                    self.register()
            elif choice == "2":
                if is_authenticated:
                    self.get_accounts()
                else:
                    self.login()
            elif is_authenticated:
                if choice == "3":
                    self.create_card()
                elif choice == "4":
                    self.get_all_cards()
                elif choice == "5":
                    self.deposit()
                elif choice == "6":
                    self.withdraw()
                elif choice == "7":
                    self.transfer()
                elif choice == "8":
                    self.get_transactions()
                elif choice == "9":
                    self.create_credit()
                elif choice == "10":
                    self.get_credits()
                elif choice == "a" and self.is_admin():
                    self.manage_users()
                elif choice == "b" and self.is_admin():
                    self.manage_roles()

if __name__ == "__main__":
    cli = BankCLI()
    cli.show_menu() 