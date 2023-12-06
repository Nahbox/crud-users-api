# CRUD users-api

Проект представляет собой HTTP сервис, который взаимодействует с базой данных PostgreSQL, реализуя CRUD API.

## Конфигурация

Конфигурация сервиса и базы данных выполняется в файле **.env**:

- ***APP_PORT*** - порт сервера
- ***POSTGRES_DB*** - название базы данных
- ***POSTGRES_USER*** - имя пользователя
- ***POSTGRES_PASSWORD*** - пароль пользователя
- ***POSTGRES_PORT*** - порт базы данных
- ***POSTGRES_HOST*** - имя хоста (при использовании docker, указывается название контейнера, в котором поднимается база данных; при локальном запуске указывается localhost)
- ***POSTGRES_MIGRATIONS_PATH*** - директория с миграциями базы данных

## Запуск

Сервис и база данных запускаются командой `docker-compose up --build`.
При первом запуске возможен случай, когда server запускается раньше postgres и не может подключиться к базе данных. Для решения данной проблемы просто перезапустите сервис:
```bash
docker-compose down
docker-compose up --build
```

## Использование сервиса

### Добавление нового пользователя:
```bash
curl -X 'POST' 'http://localhost:8080/users' -H 'Content-Type: application/json' -d '{"age": 25, "name": "John", "occupation": "Development in golang", "salary": 100000}'
curl -X 'POST' 'http://localhost:8080/users' -H 'Content-Type: application/json' -d '{"age": 50, "name": "Alex", "occupation": "Development in golang, java, c, c++ and many others.", "salary": 1000000}'
```

### Получение списка всех пользователей:
```bash
curl -X 'GET' 'http://localhost:8080/users' 
```

### Получение пользователя по его ID:
```bash
curl -X 'GET' 'http://localhost:8080/users/1'
```

### Обновление полей пользователя по его ID:
```bash
curl -X 'PUT' 'http://localhost:8080/users/1' -H 'Content-Type: application/json' -d '{"age": 20, "name": "Jake", "occupation": "Development in Java", "salary": 99999}'
```

### Удаление пользователя по его ID:
```bash
curl -X 'DELETE' 'http://localhost:8080/users/1'
```

## Swagger

Swagger с описанием API доступен после запуска сервиса по адресу:
http://localhost:8080/swagger/index.html

