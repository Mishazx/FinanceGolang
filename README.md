# FinanceGolang

Финансовый проект на Go, предоставляющий API для управления банковскими операциями.

## Описание

Этот проект представляет собой микросервисную архитектуру для обработки финансовых операций, написанную на Go. Проект включает в себя серверную часть и клиентское приложение.

## Структура проекта

```
.
├── client/              # Исходный код сервера
├── tests/client/        # Клиентское приложение
├── tests/test/          # Юнит-тесты
├── Dockerfile           # Конфигурация Docker
├── docker-compose.yml   # Конфигурациявсего проекта
├── go.mod               # Зависимости Go
└── bank_config.json     # Конфигурация банка
```

## Требования

- Go 1.21 или выше
- Docker и Docker Compose (для запуска в контейнерах)
- SQLite (для локальной разработки)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Mishazx/FinanceGolang.git
cd FinanceGolang
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите проект с помощью Docker Compose:
```bash
docker-compose up --build
```

## Разработка

Для локальной разработки:

1. Запустите сервер:
```bash
go run src/main.go
```

2. Запустите клиентское приложение:
```bash
cd tests/client
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
python3 run.py
```

