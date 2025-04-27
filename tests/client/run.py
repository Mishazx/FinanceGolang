#!/usr/bin/env python3
import os
import sys

# Добавляем путь к директории проекта в PYTHONPATH
project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.append(project_root)

from client.cli import BankCLI

if __name__ == "__main__":
    cli = BankCLI()
    cli.show_menu() 