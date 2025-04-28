import requests
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

class TransactionManager:
    def __init__(self, base_url, console, token):
        self.base_url = base_url
        self.console = console
        self.token = token

    def get_transactions(self):
        if not self.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.console.print("\n[bold blue]История транзакций[/bold blue]")
        self.console.print("=" * 50)
        
        account_id = Prompt.ask("Введите ID счета")

        response = requests.get(
            f"{self.base_url}/api/accounts/{account_id}/transactions",
            headers={"Authorization": f"Bearer {self.token}"}
        )

        if response.status_code == 200:
            transactions = response.json().get("transactions", [])
            if not transactions:
                self.console.print("[yellow]Нет транзакций по выбранному счету[/yellow]")
                return

            table = Table(show_header=True, header_style="bold magenta")
            table.add_column("Дата")
            table.add_column("Тип")
            table.add_column("Сумма")
            table.add_column("Описание")
            table.add_column("Статус")
            
            for transaction in transactions:
                # Форматируем сумму с учетом типа транзакции
                amount = float(transaction["amount"])
                formatted_amount = f"{abs(amount):.2f} ₽"
                if transaction["type"] == "withdrawal" or (transaction["type"] == "transfer" and transaction["from_account_id"] == int(account_id)):
                    formatted_amount = f"-{formatted_amount}"
                else:
                    formatted_amount = f"+{formatted_amount}"

                # Преобразуем тип транзакции для отображения
                transaction_type = transaction["type"]
                if transaction_type == "deposit":
                    display_type = "Пополнение"
                elif transaction_type == "withdrawal":
                    display_type = "Снятие"
                elif transaction_type == "transfer":
                    display_type = "Перевод"
                elif transaction_type == "credit":
                    display_type = "Оформление кредита"
                else:
                    display_type = transaction_type

                table.add_row(
                    transaction["created_at"],
                    display_type,
                    formatted_amount,
                    transaction["description"],
                    transaction["status"]
                )
            
            self.console.print(table)
        else:
            self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]") 