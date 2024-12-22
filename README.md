# Сalculator go

Данный проект является веб-сервисом с калькулятором, написанным на языке Go. Пользователь отправляет арифметическое выражение по HTTP и получает в ответ его результат.

---

## Структура проекта

- `cmd/` - директория с файлом main.go (точка входа в программу)
- `internal/application/` - директория, где находиться сервер
- `pkg/calculator/` - директория, где находиться логика самого калькулятора

---

## Установка и запуск

1. Клонируйте этот репозиторий командой:

```bash
https://github.com/MKor312/calculator_go.git
```

2. Перейдите в директорию проекта с помощью команды:

```bash
cd calculator_go/cmd
```

3. Запустите сервер командой:

```bash
go run main.go
```

---

## Использование

### Endpoint

```
POST /api/v1/calculate
```

### Заголовки

- `Content-Type: application/json`

### Тело запроса

Пример:

```json
{
  "expression": "60 / 30 - 1"
}
```

### Ответы

1. **Успешный запрос**

   **Код:** `200 OK`  
   **Пример ответа:**

   ```json
   {
     "result": "1.000000"
   }
   ```

2. **Ошибка обработки выражения**

   **Код:** `422 Unprocessable Entity`  
   **Пример ответа:**

   ```json
   {
     "invalid character in expression"
   }
   ```

3. **Некорректное тело запроса**

   **Код:** `400 Bad Request`  
   **Пример ответа:**

   ```json
   {
     "invalid expression"
   }
   ```

---

## Примеры использования

1. **Успешный запрос**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "5 + 6 * 3"
}'
```

Ответ:

```json
{
  "result": "23.000000"
}
```

2. **Ошибка: некорректное выражение**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "6 * 7y[4"
}'
```

Ответ:

```json
{
  "invalid charachter in expression"
}
```

3. **Ошибка: некорректный запрос**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10/0"
}'
```

Ответ:

```json
{
  "division by zero"
}
```
