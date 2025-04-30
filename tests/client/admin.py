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