WB level 0 task

# Сервис обработки заказов

Состоит из:
- БД Postgres
- Nats streaming server
- Http server - subscriber
- Client - publisher

Для работы требуется docker compose, taskfile, go 1.21.4

Что можно еще улучшить:
1. Тестирование БД с помощью cockroachdb testserver
2. Улучшить валидацию
3. Улучшить конфиг

