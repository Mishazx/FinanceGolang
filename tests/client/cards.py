import requests
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

class CardManager:
    def __init__(self, base_url, console, auth_manager):
        self.base_url = base_url
        self.console = console
        self.auth_manager = auth_manager

    def create_card(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Создание новой карты[/bold blue]")
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
            
            account_id = Prompt.ask("Выберите ID счета")
            
            # Создаем карту
            response = requests.post(
                f"{self.base_url}/api/cards",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"},
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

    def get_all_cards(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Список ваших карт[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            # Сначала получаем список счетов
            accounts_response = requests.get(
                f"{self.base_url}/api/accounts",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
            )
            
            if accounts_response.status_code != 200:
                self.console.print(f"[red]Ошибка при получении счетов: {accounts_response.json().get('message')}[/red]")
                return
            
            accounts = accounts_response.json().get('accounts', [])
            accounts_dict = {acc['id']: {'name': acc['name'], 'balance': acc['balance']} for acc in accounts}
            
            # Затем получаем список карт
            cards_response = requests.get(
                f"{self.base_url}/api/cards",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
            )
            
            # Отладочная информация
            self.console.print(f"[yellow]Статус ответа для карт: {cards_response.status_code}[/yellow]")
            self.console.print(f"[yellow]Ответ сервера для карт: {cards_response.text}[/yellow]")
            
            if cards_response.status_code == 200:
                data = cards_response.json()
                if not data or 'cards' not in data:
                    self.console.print("[yellow]У вас пока нет карт[/yellow]")
                    return
                    
                cards = data.get('cards', [])
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