# CRUD users-api

Проект представляет собой службу HTTP API, которая взаимодействует с базой данных PostgreSQL, реализуя операции CRUD (создание, чтение, обновление, удаление) над записями.

## Конфигурация Postgres

Конфигурация порта, на котором будет запускаться сервис, и доступа к postgres(address, user, password, database) выполняется в файле **.env**:

- ***APP_PORT*** - порт сервера
- ***POSTGRES_DB*** - название базы данных
- ***POSTGRES_USER*** - имя пользователя
- ***POSTGRES_PASSWORD*** - пароль пользователя
- ***POSTGRES_PORT*** - порт базы данных
- ***POSTGRES_HOST*** - имя хоста (при использовании docker, указывается название контейнера, на котором поднимается база данных; при локальном запуске указывается localhost)
- ***POSTGRES_MIGRATIONS_PATH*** - директория с миграциями базы данных

## Запуск

Сервис запускается командой `docker-compose up --build`.
При первом запуске сервиса возможен случай, когда server запускается раньше postgres и не может подключиться к базе данных. Для решения данной проблемы просто перезапустите сервис:
```bash
docker-compose down
docker-compose up --build
```

## Использование сервиса

### Get All Users:

Получение списка всех пользователей:
```bash
curl -X 'GET' 'http://localhost:8080/users' -H 'accept: application/json'
```

### Get User by ID:

Получение пользователя по его ID:
```bash
curl -X 'GET' 'http://localhost:8080/users/{ID}'
```

### Create User:

Добавление нового пользователя в базу данных:
```bash
curl -X 'POST' 'http://localhost:8080/users' -H 'Content-Type: application/json' -d '{"age": 0, "name": "string", "occupation": "string", "salary": 0}'
```

### Update User by ID:

Обновление полей пользователя по его ID:
```bash
curl -X 'PUT' 'http://localhost:8080/users/{ID}' -H 'Content-Type: application/json' -d '{"age": 0, "name": "string", "occupation": "string", "salary": 0}'
```

### Delete User by ID:

Удаление пользователя по его ID:
```bash
curl -X 'DELETE' 'http://localhost:8080/users/{ID}'
```

## Тесты

Тесты находятся в директории **/internal/config**.

Для запуска **config_test.go** введите:

    go test -v -run FromEnv

Для запуска **postgres_test.go** введите:

    go test -v -run PgConfig

> Запуск тестов производится непосредственно из директории, в которой они находятся.
