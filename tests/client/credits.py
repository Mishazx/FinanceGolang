import requests
from rich.console import Console
from rich.prompt import Prompt, Confirm
from rich.table import Table

class CreditManager:
    def __init__(self, base_url, console, auth_manager):
        self.base_url = base_url
        self.console = console
        self.auth_manager = auth_manager

    def create_credit(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Оформление кредита[/bold blue]")
        self.console.print("=" * 50)
        
        # Получаем список счетов
        try:
            response = requests.get(
                f"{self.base_url}/api/accounts",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
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
            description = str(Prompt.ask("Введите описание кредита", default="Потребительский кредит"))
            
            # Создаем JSON для отправки
            json_data = {
                "account_id": int(account_id),
                "amount": amount,
                "term_months": term_months,
                "description": description
            }
            
            # Выводим информацию о JSON для отладки
            self.console.print(f"[bold]Отправляемые данные:[/bold] {json_data}")
            
            try:
                response = requests.post(
                    f"{self.base_url}/api/credits",
                    headers={"Authorization": f"Bearer {self.auth_manager.token}"},
                    json=json_data
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
                    # Отображаем детальную информацию об ошибке
                    error_text = response.text
                    self.console.print(f"[red]Статус код: {response.status_code}[/red]")
                    self.console.print(f"[red]Ответ сервера: {error_text}[/red]")
                    
            except Exception as e:
                self.console.print(f"[red]Ошибка: {str(e)}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")

    def get_credits(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Мои кредиты[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            response = requests.get(
                f"{self.base_url}/api/credits",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
            )
            
            if response.status_code == 200:
                data = response.json()
                if not data or 'credits' not in data:
                    self.console.print("[yellow]У вас пока нет кредитов[/yellow]")
                    return
                    
                credits = data.get('credits', [])
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
                self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")

    def get_payment_schedule(self, credit_id):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print(f"\n[bold blue]График платежей по кредиту #{credit_id}[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            response = requests.get(
                f"{self.base_url}/api/credits/{credit_id}/schedule",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
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
                        payment['due_date'],
                        f"{payment['principal']:.2f} ₽",
                        f"{payment['interest']:.2f} ₽",
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

    def process_payment(self, credit_id, payment_number):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print(f"\n[bold blue]Оплата платежа #{payment_number} по кредиту #{credit_id}[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            response = requests.post(
                f"{self.base_url}/api/credits/{credit_id}/payment",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"},
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