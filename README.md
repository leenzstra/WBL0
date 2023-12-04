WB level 0 task

# Сервис обработки заказов

Состоит из:
- БД Postgres
- Nats streaming server
- Http server / subscriber / UI
- Client / publisher

Запуск через docker-compose, taskfile, go 1.21.4

API - localhost/api/v1/orders/:uid

UI - localhost/ui

- cmd/sub/server.go - сервис
- cmd/pub/publisher.go - клиент

Что можно еще сделать лучше:
1. Тестировать БД с помощью cockroachdb testserver
2. Улучшить валидацию
3. Делать коммиты более атомарными

