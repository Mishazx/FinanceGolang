import requests
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

class AdminManager:
    def __init__(self, base_url, console, auth_manager):
        self.base_url = base_url
        self.console = console
        self.auth_manager = auth_manager

    def manage_users(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Управление пользователями[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            response = requests.get(
                f"{self.base_url}/api/admin/users",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
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
                        headers={"Authorization": f"Bearer {self.auth_manager.token}"},
                        json={"role": role}
                    )
                else:
                    response = requests.delete(
                        f"{self.base_url}/api/admin/users/{user_id}/roles/{role}",
                        headers={"Authorization": f"Bearer {self.auth_manager.token}"}
                    )
                
                if response.status_code in [200, 201]:
                    self.console.print("[green]Операция выполнена успешно![/green]")
                else:
                    self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")

    def manage_roles(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Управление ролями[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            response = requests.get(
                f"{self.base_url}/api/admin/roles",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
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
                        headers={"Authorization": f"Bearer {self.auth_manager.token}"},
                        json={
                            "name": name,
                            "description": description
                        }
                    )
                else:
                    role_name = Prompt.ask("Введите название роли для удаления")
                    response = requests.delete(
                        f"{self.base_url}/api/admin/roles/{role_name}",
                        headers={"Authorization": f"Bearer {self.auth_manager.token}"}
                    )
                
                if response.status_code in [200, 201]:
                    self.console.print("[green]Операция выполнена успешно![/green]")
                else:
                    self.console.print(f"[red]Ошибка: {response.json().get('message')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]")

    def test_scheduler(self):
        if not self.auth_manager.token:
            self.console.print("[red]Ошибка: Необходима авторизация[/red]")
            return
        
        self.console.print("\n[bold blue]Тестирование шедулера[/bold blue]")
        self.console.print("=" * 50)
        
        try:
            # Получаем список всех кредитов
            response = requests.get(
                f"{self.base_url}/api/admin/credits",
                headers={"Authorization": f"Bearer {self.auth_manager.token}"}
            )
            
            if response.status_code == 200:
                credits = response.json().get('credits', [])
                
                if not credits:
                    self.console.print("[yellow]Нет активных кредитов для проверки[/yellow]")
                    return
                
                # Показываем список кредитов
                table = Table(show_header=True, header_style="bold magenta")
                table.add_column("ID")
                table.add_column("Пользователь")
                table.add_column("Сумма")
                table.add_column("Статус")
                table.add_column("Дата начала")
                table.add_column("Дата окончания")
                
                for credit in credits:
                    table.add_row(
                        str(credit['id']),
                        str(credit['user_id']),
                        f"{credit['amount']:.2f} ₽",
                        credit['status'],
                        credit['start_date'],
                        credit['end_date']
                    )
                
                self.console.print(table)
                
                # Запускаем проверку платежей
                self.console.print("\n[bold]Запуск проверки платежей...[/bold]")
                response = requests.post(
                    f"{self.base_url}/api/admin/scheduler/check-payments",
                    headers={"Authorization": f"Bearer {self.auth_manager.token}"}
                )
                
                if response.status_code == 200:
                    result = response.json()
                    self.console.print("[green]Проверка платежей выполнена![/green]")
                    
                    # Показываем результаты
                    result_table = Table(show_header=True, header_style="bold green")
                    result_table.add_column("Метрика")
                    result_table.add_column("Значение")
                    
                    for key, value in result.items():
                        result_table.add_row(key, str(value))
                    
                    self.console.print(result_table)
                else:
                    self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            else:
                self.console.print(f"[red]Ошибка: {response.json().get('error')}[/red]")
            
        except Exception as e:
            self.console.print(f"[red]Ошибка: {str(e)}[/red]") 