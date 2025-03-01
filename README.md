# Сalculator go

Данный проект является простым HTTP сервером с калькулятором, написанным на языке Go. Калькулятор принимает арифметическое выражение по HTTP и возвращает его результат. Если присутствует ошибка, то калькулятор возвращает эту ошибку.

---

## Структура проекта

- `cmd/` - директория с файлом main.go (точка входа в программу)
- `internal/application/` - директория, где находиться сервер
- `pkg/calc/` - директория, где находиться логика самого калькулятора

---

## Установка и запуск

1. Клонируйте этот репозиторий командой:

```bash
git clone https://github.com/MKode312/calculator_go.git
```

2. Перейдите в директорию проекта с помощью команды:

```bash
cd calculator_go/cmd
```

3. Для запуска сервера используйте команду:

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
     "result": 1
   }
   ```

2. **Ошибка обработки выражения**

   **Код:** `422 Unprocessable Entity`  
   **Пример ответа:**

   ```json
   {
     "error": "Expression is not valid"
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
  "result": 23
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
  "error": "Expression is not valid"
}
```

### Примечание

При возникновении непредвиденных ошибок со стороны сервера, веб-сервис вернёт ошибку:

```json
{
  "error": "Internal server error"
}
``` 
**Код:** `500 Internal Server Error`
