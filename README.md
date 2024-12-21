# Калькулятор Web-сервис

Веб-сервис для вычисления арифметических выражений через HTTP API.

## Возможности

- Поддерживает базовые арифметические операции (+, -, *, /)
- Работает с десятичными числами
- Поддерживает скобки для группировки операций
- Возвращает результаты в формате JSON

## API

### Вычисление выражения

**Endpoint:** `POST /api/v1/calculate`

**Тело запроса:**

```json
{
"expression": "2 + 2 * 2"
}
```

**Успешный ответ (200 OK):**

```json
{
"result": 6
}
```

**Ошибка валидации (422 Unprocessable Entity):**

```json
{
"error": "Expression is not valid"
}
```

### Возможные сообщения об ошибках

- `"Expression cannot be empty"` - пустое выражение
- `"Division by zero is not allowed"` - деление на ноль
- `"Consecutive operators are not allowed"` - два оператора подряд
- `"Parentheses are mismatched"` - непарные скобки
- `"Expression is not valid"` - некорректное выражение
- `"Invalid request body"` - неверный формат JSON
- `"Internal server error"` - внутренняя ошибка сервера

## Установка и запуск

1. Клонировать репозиторий:

```bash
git clone https://github.com/hxzzf/calc_go.git
```

2. Перейти в директорию проекта:

```bash
cd calc_go
```

3. Запустить сервис:

```bash
go run cmd/main.go
```

По умолчанию сервис запускается на порту 8080. Чтобы изменить порт, установите переменную окружения PORT:

### Linux/MacOS:
```bash
PORT=3000 go run cmd/main.go
```

### Windows PowerShell:
```powershell
$env:PORT=3000; go run cmd/main.go
```

### Windows CMD:
```cmd
set PORT=3000
go run cmd/main.go
```
Примечание: в CMD необходимо выполнить команды по отдельности.

## Остановка сервера

- Остановка: нажмите `Ctrl+C` для остановки сервера

## Примеры использования

### Windows (cmd)
```bash
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"2 + 2 * 2\"}"
```

### Linux/MacOS
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "2 + 2 * 2"}'
```

Другие примеры выражений:

```bash
# Выражение со скобками
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"(2 + 3) * 4\"}"

# Деление на ноль
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"1 / 0\"}"

# Некорректное выражение
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"2 + + 2\"}"

# Ошибка «Что-то пошло не так» (Internal Server Error)
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"1,7976931348623157 * 2\"}"
```

Ожида��мые ответы:
```json
// Успешное вычисление
{"result": 20}

// Ошибка деления на ноль
{"error": "Division by zero is not allowed"}

// Ошибка некорректного выражения
{"error": "Consecutive operators are not allowed"}

// Ошибка «Что-то пошло не так»
{"error": "Internal server error"}
```

## Разработка и тестирование

Запуск тестов:

```bash
go test ./...
```

