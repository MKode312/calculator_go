# Calculator Go

Данный проект является простым HTTP сервером с калькулятором, написанным на языке Go. Калькулятор принимает арифметическое выражение по HTTP и возвращает его результат. Если присутствует ошибка, то калькулятор возвращает эту ошибку.

## Структура проекта

- `cmd/` - директория с файлом main.go (точка входа в программу)
- `internal/application/` - директория, где находиться сервер
- `pkg/calculator/` - директория, где находиться логика самого калькулятора

---

## Установка

Для того чтобы установить и запустить проект, выполните следующие шаги:

1. Клонируйте этот репозиторий командой: `git clone https://github.com/MKor312/calculator_go.git`
2. Перейдите в директорию проекта c помощью команды: `cd calculator_go/cmd`
3. Для запуска калькулятора выполните следующую команду: `go run main.go`

---

## Использование

Вы можете отдавать серверу запросы с помощью curl, вот несколько примеров запросов для Windows (Command Prompt): 

1. `curl --location "http://localhost:8080/api/v1/calculate" ^
--header "Content-Type: application/json" ^
--data "{ \"expression\": \"2 + 2 * 4 \" }"`	

2. `curl --location "http://localhost:8080/api/v1/calculate" ^
--header "Content-Type: application/json" ^
--data "{ \"expression\": \"8 + 679\" }"`

3. `curl --location "http://localhost:8080/api/v1/calculate" ^
--header "Content-Type: application/json" ^
--data "{ \"expression\": \"(4 + 5) * 6\" }"`

### Для macOS и Linux:

1. `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "156 - 36" }'`

2. `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "9 * 8" }'`

3. `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "(65 + 35) / 100" }'`

---

## Примеры

Вот некоторые примеры использования для Windows (Command Prompt) с ожидаемыми результатами:

- Ввод: `curl --location "http://localhost:8080/api/v1/calculate" ^ 
--header "Content-Type: application/json" ^ 
--data "{ \"expression\": \"(8 + 6) * 4\" }"`
  - Вывод: `result: 56.000000 `
- Код [200]
  
- Ввод: `curl --location "http://localhost:8080/api/v1/calculate" ^
 --header "Content-Type: application/json" ^
 --data "{ \"expression\": \"10/0\" }"`
  - Вывод: `division by zero`
- Код [400]
  
- Ввод: `curl --location "http://localhost:8080/api/v1/calculate" ^ 
--header "Content-Type: application/json" ^ 
--data "{ \"expression\": \"8 + abc\" }"`
  - Вывод: `invalid character in expression`
- Код [422]

### Для macOS и Linux:

- Ввод: `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "(65 + 35) / 100" }'`
  - Вывод: `result: 1.000000 `
- Код [200]
  
- Ввод: `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "7**6" }'`
  - Вывод: `two consecutive operators`
- Код [400]
  
- Ввод: `curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{ "expression": "19 - n" }'`
  - Вывод: `invalid character in expression`
- Код [422]

В случае возникновения непредвиденных ошибок со стороны сервера, вывод будет таков: 
`Internal server error`
- Код: [500]


 
