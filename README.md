Go Pet Shop
Многофункциональный сервис для управления интернет-магазином зоотоваров. Проект реализован на Go с использованием современных технологий: PostgreSQL, ClickHouse, Kafka, Docker, Chi, Slog и других.

🚀 Возможности
- Хранение и обработка заказов, пользователей, транзакций, истории покупок
- Высокопроизводительная аналитика через ClickHouse
- Асинхронная обработка событий через Kafka
- REST API для взаимодействия с фронтендом/мобильными приложениями
- Миграции и автоматизация через Taskfile
- Гибкая конфигурация через YAML
🛠️ Технологии
- Go 1.21+
- PostgreSQL (основная БД)
- ClickHouse (аналитика)
- Kafka (очереди событий)
- Chi (HTTP роутер)
- Slog (логирование)
- Docker (контейнеризация)
- golang-migrate (миграции)
- pgxpool (pool для PostgreSQL)
- go-chi/render (JSON-ответы)
- godotenv, cleanenv (конфиг)
- Taskfile (автоматизация)
  
📦 Архитектура

```
cmd/
  app/         # Основной HTTP сервер
  migrator/    # CLI для миграций
internal/
  config/      # Загрузка конфигов
  delivery/    # HTTP хендлеры
  domain/      # Модели данных
  lib/logger/  # Логирование
  storage/     # Репозитории и работа с БД
config/
  local.yaml   # Конфиг подключения
migrations/    # SQL миграции
Taskfile.yaml  # Автоматизация задач
```

⚡ Быстрый старт

1.Клонируйте репозиторий:
git clone https://github.com/flizity/go_pet_shop_v2.git
cd go_pet_shop_v2

2.Запустите необходимые сервисы через Docker:
docker-compose -f docker-compose.kafka.yaml up -d

3.Настройте конфиг local.yaml под свои параметры.

4.Примените миграции:
go run cmd/migrator/main.go up

5.Запустите сервер:
go run cmd/app/main.go

📊 Аналитика и события
- ClickHouse используется для хранения истории заказов и аналитики.
- Kafka — для асинхронной обработки событий (например, заказов).

📝 Миграции
- Все миграции хранятся в папке migrations.
- Используется golang-migrate для применения изменений схемы.

🏗️ Автоматизация
- Taskfile.yaml позволяет запускать типовые задачи (миграции, тесты, сборка).

💡 Контакты и поддержка
- Автор: flizity
- Вопросы и предложения — через Issues на GitHub.
