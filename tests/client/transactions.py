import requests
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

class TransactionManager:
    def __init__(self, base_url, console, auth_manager):
        self.base_url = base_url
        self.console = console
        self.auth_manager = auth_manager

    def get_transactions(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return

        self.console.print("\n[bold blue]История транзакций[/bold blue]")
        self.console.print("=" * 50)
        
        account_id = Prompt.ask("Введите ID счета")

        response = requests.get(
            f"{self.base_url}/api/accounts/{account_id}/transactions",
            headers={"Authorization": f"Bearer {self.auth_manager.token}"}
        )

        if response.status_code == 200:
            transactions = response.json().get("transactions", [])
            if not transactions:
                self.console.print("[yellow]Нет транзакций по выбранному счету[/yellow]")
                return

            # Отладочный вывод
            self.console.print("\n[bold yellow]Отладочная информация:[/bold yellow]")
            for t in transactions:
                self.console.print(f"Транзакция: {t}")

            table = Table(show_header=True, header_style="bold magenta")
            table.add_column("Дата")
            table.add_column("Тип")
            table.add_column("Сумма")
            table.add_column("Описание")
            table.add_column("Статус")
            
            for transaction in transactions:
                # Получаем сумму как число
                amount = float(transaction["amount"])
                
                # Определяем знак суммы в зависимости от типа транзакции и роли счета
                if transaction["type"] == "TRANSFER":
                    # Для переводов проверяем, является ли счет отправителем или получателем
                    from_id = transaction.get("from_account_id")
                    to_id = transaction.get("to_account_id")
                    self.console.print(f"\n[bold yellow]Отладка перевода:[/bold yellow]")
                    self.console.print(f"from_id: {from_id}, to_id: {to_id}, account_id: {account_id}")
                    
                    if from_id and int(from_id) == int(account_id):
                        # Если счет отправитель - сумма отрицательная
                        formatted_amount = f"-{abs(amount):.2f} ₽"
                    elif to_id and int(to_id) == int(account_id):
                        # Если счет получатель - сумма положительная
                        formatted_amount = f"+{abs(amount):.2f} ₽"
                    else:
                        # Если счет не участвует в переводе (не должно происходить)
                        formatted_amount = f"{abs(amount):.2f} ₽"
                elif transaction["type"] == "WITHDRAWAL":
                    # Для снятий сумма всегда отрицательная
                    formatted_amount = f"-{abs(amount):.2f} ₽"
                elif transaction["type"] == "PAYMENT":
                    # Для оформления кредита сумма всегда положительная
                    formatted_amount = f"-{abs(amount):.2f} ₽"
                else:
                    # Для пополнений сумма всегда положительная
                    formatted_amount = f"+{abs(amount):.2f} ₽"

                # Преобразуем тип транзакции для отображения
                transaction_type = transaction["type"]
                if transaction_type == "DEPOSIT":
                    display_type = "Пополнение"
                elif transaction_type == "WITHDRAWAL":
                    display_type = "Снятие"
                elif transaction_type == "TRANSFER":
                    display_type = "Перевод"
                elif transaction_type == "CREDIT":
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