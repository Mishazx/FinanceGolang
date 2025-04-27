import sys
import requests
from rich.console import Console
from rich.prompt import Prompt, Confirm
from auth import AuthManager
from accounts import AccountManager
from cards import CardManager
from credits import CreditManager
from transactions import TransactionManager
from admin import AdminManager

class BankCLI:
    def __init__(self):
        self.base_url = "http://localhost:8080"
        self.console = Console()
        self.auth_manager = AuthManager(self.base_url, self.console)
        
        # Отладочная информация
        self.console.print(f"[yellow]Токен при инициализации: {self.auth_manager.token}[/yellow]")
        
        self.account_manager = AccountManager(self.base_url, self.console, self.auth_manager.token)
        self.card_manager = CardManager(self.base_url, self.console, self.auth_manager.token)
        self.credit_manager = CreditManager(self.base_url, self.console, self.auth_manager.token)
        self.transaction_manager = TransactionManager(self.base_url, self.console, self.auth_manager.token)
        self.admin_manager = AdminManager(self.base_url, self.console, self.auth_manager.token)

    def clear_screen(self):
        self.console.clear()

    def print_header(self, title):
        self.clear_screen()
        self.console.print(f"\n[bold blue]{title}[/bold blue]")
        self.console.print("=" * 50)

    def show_menu(self):
        while True:
            # Очищаем экран и показываем заголовок
            self.clear_screen()
            self.console.print("\n[bold blue]Банковский терминал[/bold blue]")
            self.console.print("=" * 50)
            
            is_authenticated = self.auth_manager.check_auth()
            
            if is_authenticated:
                role_text = ", ".join(self.auth_manager.user_roles) if self.auth_manager.user_roles else "нет ролей"
                self.console.print(f"[green]Вы авторизованы[/green] (Роли: {role_text})")
            else:
                self.console.print("[yellow]Вы не авторизованы[/yellow]")
                self.auth_manager.token = None
                self.auth_manager.user_roles = []
                self.auth_manager.save_config()
            
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
                
                if self.auth_manager.is_admin():
                    self.console.print("a. Управление пользователями")
                    self.console.print("b. Управление ролями")
                
                self.console.print("l. Выйти из аккаунта")
                self.console.print("0. Завершить программу")
            
            if not is_authenticated:
                choices = ["0", "1", "2"]
            else:
                choices = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "l"]
                if self.auth_manager.is_admin():
                    choices.extend(["a", "b"])
            
            choice = Prompt.ask("\nВыберите действие", choices=choices)
            
            if choice == "0":
                if Confirm.ask("Вы уверены, что хотите завершить программу?"):
                    sys.exit()
            elif choice == "l" and is_authenticated:
                self.auth_manager.logout()
            elif choice == "1":
                if is_authenticated:
                    self.account_manager.create_account()
                else:
                    self.auth_manager.register()
            elif choice == "2":
                if is_authenticated:
                    self.account_manager.get_accounts()
                else:
                    self.auth_manager.login()
            elif is_authenticated:
                if choice == "3":
                    self.card_manager.create_card()
                elif choice == "4":
                    self.card_manager.get_all_cards()
                elif choice == "5":
                    self.account_manager.deposit()
                elif choice == "6":
                    self.account_manager.withdraw()
                elif choice == "7":
                    self.account_manager.transfer()
                elif choice == "8":
                    self.transaction_manager.get_transactions()
                elif choice == "9":
                    self.credit_manager.create_credit()
                elif choice == "10":
                    self.credit_manager.get_credits()
                elif choice == "a" and self.auth_manager.is_admin():
                    self.admin_manager.manage_users()
                elif choice == "b" and self.auth_manager.is_admin():
                    self.admin_manager.manage_roles()
            
            # Добавляем паузу после выполнения команды
            input("\nНажмите Enter для продолжения...")

if __name__ == "__main__":
    cli = BankCLI()
    cli.show_menu() 