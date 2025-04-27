import requests
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

class AccountManager:
    def __init__(self, base_url, console, token):
        self.base_url = base_url
        self.console = console
        self.token = token

    def create_account(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.console.print("\n[bold blue]Создание счета[/bold blue]")
        self.console.print("=" * 50)
        
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

    def get_accounts(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.console.print("\n[bold blue]Мои счета[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            # Отладочная информация о токене
            self.console.print(f"[yellow]Используемый токен: {self.token}[/yellow]")
            
            url = f"{self.base_url}/api/accounts"
            self.console.print(f"[yellow]Запрос к URL: {url}[/yellow]")
            
            headers = {
                "Authorization": f"Bearer {self.token}",
                "Content-Type": "application/json"
            }
            self.console.print(f"[yellow]Заголовки запроса: {headers}[/yellow]")
            
            response = requests.get(url, headers=headers)
            
            # Отладочная информация
            self.console.print(f"[yellow]Статус ответа: {response.status_code}[/yellow]")
            self.console.print(f"[yellow]Ответ сервера: {response.text}[/yellow]")
            
            if response.status_code == 200:
                data = response.json()
                if not data or 'accounts' not in data:
                    self.console.print("[yellow]У вас пока нет счетов[/yellow]")
                    return
                    
                accounts = data.get('accounts', [])
                if not accounts:
                    self.console.print("[yellow]У вас пока нет счетов[/yellow]")
                    return

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
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")
            import traceback
            self.console.print(f"[red]Трассировка: {traceback.format_exc()}[/red]")

    def deposit(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.console.print("\n[bold blue]Пополнение счета[/bold blue]")
        self.console.print("=" * 50)
        
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

        self.console.print("\n[bold blue]Снятие средств[/bold blue]")
        self.console.print("=" * 50)
        
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

        self.console.print("\n[bold blue]Перевод средств[/bold blue]")
        self.console.print("=" * 50)
        
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