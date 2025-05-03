# FinanceGolang - Банковское приложение

Банковское приложение на Go с поддержкой кредитов, транзакций и управления счетами.

## Функциональность

- 💰 Управление счетами
- 💳 Кредиты и платежи по кредитам
- 📊 История транзакций
- 🔒 JWT аутентификация
- 📱 REST API

## Технологический стек

- Go 1.21+ - основной язык разработки
- PostgreSQL - реляционная база данных
- GORM - ORM для работы с базой данных
- Gin Web Framework - веб-фреймворк
- JWT - аутентификация и авторизация
- Docker - контейнеризация
- Python - тестовый клиент

## Варианты запуска

### 1. Запуск в Docker (рекомендуемый способ)

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Mishazx/FinanceGolang.git
cd FinanceGolang
```

2. Поменять переменные если нужно

3. Запустите приложение и базу данных:
```bash
docker-compose up --build
```

Приложение будет доступно по адресу: `http://localhost:8080`

### 2. Локальный запуск с PostgreSQL в Docker

1. Установите зависимости:
```bash
go mod download
```

2. Запустите только PostgreSQL:
```bash
docker-compose up -d postgres
```

3. Создайте файл `.env`:
```env
# Настройки сервера
SERVER_HOST=localhost
SERVER_PORT=8080

# База данных
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bank
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24

# Приложение
APP_ENV=development
APP_DEBUG=true
```

4. Запустите приложение:
```bash
go run src/main.go
```

### 3. Полностью локальный запуск

1. Установите зависимости:
```bash
go mod download
```


2. Запустите приложение:
```bash
go run src/main.go
```

## Тестирование API

### Через Python клиент

1. Установите зависимости клиента:
```bash
cd tests/client
python3 -m venv .venv
source .venv/bin/activate  # для Linux/Mac
# или
.venv\Scripts\activate     # для Windows
pip install -r requirements.txt
```

2. Запустите клиент:
```bash
python main.py
```

### Через curl

1. Регистрация нового пользователя:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

2. Вход и получение токена:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

3. Использование полученного токена:
```bash
# Сохраните токен из предыдущего ответа
export TOKEN="полученный_токен"

# Пример запроса с токеном
curl -X GET http://localhost:8080/api/accounts \
  -H "Authorization: Bearer $TOKEN"
```

## Конфигурация

Приложение настраивается через переменные окружения или файл `.env`:

### Основные настройки

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| SERVER_HOST | Хост сервера | localhost |
| SERVER_PORT | Порт сервера | 8080 |
| DB_TYPE | Тип базы данных | postgres |
| DB_HOST | Хост базы данных | localhost |
| DB_PORT | Порт базы данных | 5432 |
| DB_USER | Пользователь БД | postgres |
| DB_PASSWORD | Пароль БД | postgres |
| DB_NAME | Имя базы данных | bank |
| DB_SSLMODE | Режим SSL | disable |
| JWT_SECRET | Секретный ключ JWT | your-secret-key |
| JWT_EXPIRATION | Время жизни токена (часы) | 24 |
| APP_ENV | Окружение (development/production) | development |
| APP_DEBUG | Режим отладки | true |

## API Endpoints

### Аутентификация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `POST /api/auth/refresh` - Обновление токена

### Счета
- `GET /api/accounts` - Список счетов
- `POST /api/accounts` - Создание счета
- `GET /api/accounts/:id` - Информация о счете
- `GET /api/accounts/:id/transactions` - История транзакций

### Кредиты
- `POST /api/credits` - Оформление кредита
- `GET /api/credits` - Список кредитов
- `GET /api/credits/:id` - Информация о кредите
- `POST /api/credits/:id/payment` - Внесение платежа
- `GET /api/credits/:id/schedule` - График платежей

### Транзакции
- `POST /api/transactions` - Создание транзакции
- `GET /api/transactions` - История транзакций
- `GET /api/transactions/:id` - Детали транзакции

## Безопасность

- Все API-endpoints защищены JWT аутентификацией (кроме регистрации и входа)
- Пароли хешируются перед сохранением
- SSL режим для базы данных настраивается через конфигурацию
- Валидация всех входящих данных

## Структура проекта
```
.
├── src/
│   ├── config/            # Конфигурация приложения
│   ├── controller/        # HTTP контроллеры
│   ├── database/          # Конфигурация и миграции БД
│   ├── dto/               # Data Transfer Objects
│   ├── model/             # Модели данных
│   ├── repository/        # Слой доступа к данным
│   ├── security/          # Аутентификация и авторизация
│   └── service/           # Бизнес-логика
├── tests/
│   ├── client/            # Тестовый клиент на Python
│   └── test/              # Интеграционные тесты
├── .template.env          # Env шаблон
├── docker-compose.yml     # Docker конфигурация
├── Dockerfile            # Сборка приложения
└── README.md             # Документация
```
