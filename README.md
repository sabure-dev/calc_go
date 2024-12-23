# Calculator API 🧮

Простой и эффективный REST API калькулятор. Поддерживает базовые математические операции и работу со скобками.

## Возможности 🚀

- Базовые математические операции:
  - Сложение (`+`)
  - Вычитание (`-`)
  - Умножение (`*`)
  - Деление (`/`)
- Поддержка скобок для приоритета операций
- Валидация выражений
- Подробные сообщения об ошибках
- Поддержка docker и docker-compose
- Конфигурация через файл `.env`
- Логирование запросов
- CORS поддержка
- Swagger UI
- Тестирование

## Установка и запуск 🛠️

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/sabure-dev/calc_go.git
   cd calc_go
   ```
2. Создайте файл `.env` в корневой директории, чтобы задать порт для сервера. По умолчанию используется порт 8080:
   ```env
   PORT=8080
   ```
   
### Локальный запуск
   Запуск приложения:
   ```bash
   go run cmd/main.go
   ```

### Docker (предпочтительный способ запуска)🐳
   Сборка и запуск через docker compose:
   ```bash
   docker compose up --build
   ```

   Или через Docker напрямую:
   ```bash
   docker build -t calc-api .
   docker run -p 8080:8080 calc-api
   ```


## Использование 📝

### Endpoint

`POST /api/v1/calculate`

### Формат запроса

```json
{
  "expression": "2+2*2"
}
```

### Формат ответа

Успешный ответ:
**Successful Response (200):**
```json
{
  "result": 6
}
```

Ответы с ошибкой:

**Validation Error (422):**
```json
{
  "error": "Expression is not valid"
}
```

**Method Not Allowed (405):**
```json
{
  "error": "Method not allowed"
}
```

**Internal Server Error (500):**
```json
{
  "error": "Internal server error"
}
```

### Примеры запросов

Простое выражение
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "2+2"}'
```
Выражение со скобками
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "(2+2)*2"}'
```

## API документация 📚

API документация доступна через Swagger UI после запуска проекта:
- Docker: http://localhost:8081

### Swagger UI

Swagger UI предоставляет интерактивный интерфейс для:
- Просмотра всех доступных эндпоинтов
- Тестирования API запросов
- Изучения форматов запросов и ответов
- Просмотра возможных кодов ответов и ошибок

## Обработка ошибок 🚨

API возвращает понятные сообщения об ошибках в следующих случаях:

- Некорректный метод запроса (только POST)
- Некорректный JSON в теле запроса
- Выражение слишком короткое (минимум 3 символа)
- Некорректное расположение операторов
- Последовательные операторы (например, "2++2")
- Деление на ноль
- Некорректные символы в выражении
- Ошибки в выражениях внутри скобок

## Разработка 👨‍💻

### Структура проекта
    ├── cmd/
    │ └── main.go # Точка входа
    ├── internal/
    │ └── application/ # Внутренняя логика приложения
    │ ├──── application.go
    │ ├──── handlers.go
    │ └──── middleware.go
    ├── pkg/
    │ └── calculation/ # Логика вычислений
    │ ├──── calculation.go
    │ └──── errors.go
    ├── Dockerfile
    ├── docker-compose.yml
    ├── go.mod
    ├── .gitignore
    ├── swagger.yaml # OpenAPI спецификация
    ├── LICENSE
    ├── go.sum
    └── README.md

### Запуск тестов

```bash
go test ./...
```

## Алгоритм вычислений 🧮

Калькулятор использует следующий алгоритм для вычисления выражений:

1. **Предварительная валидация**:
   - Проверка минимальной длины выражения (3 символа)
   - Проверка корректности расположения операторов
   - Проверка на последовательные операторы

2. **Обработка скобок**:
   - Поиск самых вложенных скобок
   - Рекурсивное вычисление выражений внутри скобок
   - Замена скобочного выражения на результат

3. **Соблюдение приоритета операций**:
   - Первый проход: вычисление умножения и деления (слева направо)
   - Второй проход: вычисление сложения и вычитания (слева направо)

4. **Обработка чисел и операторов**:
   - Парсинг чисел с учетом отрицательных значений
   - Последовательное применение операторов
   - Проверка деления на ноль

### Пример работы алгоритма

Для выражения `(2+2)*2`:

1. Находит выражение в скобках `(2+2)`
2. Рекурсивно вычисляет `2+2 = 4`
3. Заменяет `(2+2)` на `4`
4. Получает выражение `4*2`
5. Выполняет умножение: `4*2 = 8`
6. Возвращает результат: `8`

Основная логика вычислений реализована в файле: `pkg/calculation/calculation.go`.

## Лицензия 📄

MIT License - см. файл [LICENSE](LICENSE)

## Вклад в проект 🤝

Приветствуются любые предложения по улучшению проекта! Для этого:

1. Форкните репозиторий
2. Создайте ветку для новой функциональности
3. Внесите изменения
4. Отправьте pull request

## Контакты 📧

- GitHub: [@sabure-dev](https://github.com/sabure-dev)
