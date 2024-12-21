# Calculator Go

Данный проект является простым HTTP сервером с калькулятором, написанным на языке Go. Калькулятор принимает арифметическое выражение по HTTP и возвращает его результат.

## Установка

Для того чтобы установить и запустить проект, выполните следующие шаги:

1. Клонируйте этот репозиторий командой: `git clone https://github.com/MKor312/calculator_go.git`
2. Перейдите в директорию проекта c помощью команды: `cd calculator_go/cmd`
3. Для запуска калькулятора выполните следующую команду: `go run main.go`

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

Также немного примеров запросов для macOS и Linux:

1. `curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2 + 2 * 2"
}'`

2. `curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(78 + 83) / 160"
}'`

3. `curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "7 * 62"
}'`

## Примеры

Вот некоторые примеры использования с ожидаемыми результатами:

- Ввод: `3 + 4`
  - Вывод: `7`
- Код: 200
  
- Ввод: `10/0`
  - Вывод: `division by zero`
- Код: 400
  
- Ввод: `8 + a`
  - Вывод: `invalid character in expression`
- Код: 422
  
