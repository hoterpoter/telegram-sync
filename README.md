# telegram-sync

Проект для синхронизации сообщений и фотографий из Telegram в GitHub репозиторий.

## Описание

Это приложение на Go с использованием фреймворка Fiber позволяет автоматически сохранять сообщения и фотографии из Telegram в GitHub репозиторий. Бот принимает webhook'и от Telegram и сохраняет контент в структурированном виде.

## Возможности

- 📝 Синхронизация текстовых сообщений в GitHub
- 📸 Синхронизация фотографий в GitHub
- 🗂️ Организованное хранение по датам и чатам
- 🔄 Автоматическая обработка webhook'ов
- ✅ Уведомления об успешном сохранении
- 🚀 Быстрая работа благодаря Fiber

## Требования

- Go 1.21 или выше
- Telegram Bot Token (получить у [@BotFather](https://t.me/botfather))
- GitHub Personal Access Token с правами на запись в репозиторий
- Публичный URL для webhook'ов (можно использовать ngrok для разработки)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/hoterpoter/telegram-sync.git
cd telegram-sync
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

4. Заполните `.env` файл:
```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook/telegram
GITHUB_TOKEN=your_github_token
GITHUB_OWNER=your_username
GITHUB_REPO=your_repo
GITHUB_BRANCH=main
PORT=3000
```

## Запуск

### Использование Makefile (рекомендуется)

```bash
# Посмотреть все доступные команды
make help

# Собрать приложение
make build

# Запустить приложение
make run

# Форматировать и проверить код
make lint
```

### Локальная разработка

```bash
# Загрузите переменные окружения
export $(cat .env | xargs)

# Запустите приложение
go run main.go
```

### Сборка

```bash
go build -o telegram-sync
./telegram-sync
```

### Docker

```bash
# Сборка Docker образа
docker build -t telegram-sync:latest .

# Запуск с Docker Compose
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

Или с использованием Makefile:
```bash
make docker-build
make docker-run
make docker-logs
make docker-stop
```

## Настройка Telegram Webhook

После запуска приложения, настройте webhook для вашего бота:

```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=<YOUR_WEBHOOK_URL>"
```

Или через API:
```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-domain.com/webhook/telegram"}'
```

## Структура сохранения

Сообщения и фотографии сохраняются в следующей структуре:

```
messages/
  chat_{chat_id}/
    {date}/
      msg_{message_id}.txt

photos/
  chat_{chat_id}/
    {date}/
      photo_{message_id}.jpg
      photo_{message_id}_caption.txt
```

## Использование

1. Добавьте бота в чат или отправьте ему личное сообщение
2. Отправьте текстовое сообщение или фотографию
3. Бот автоматически сохранит контент в GitHub
4. Вы получите подтверждение об успешном сохранении

## API Endpoints

- `GET /` - Проверка статуса сервиса
- `POST /webhook/telegram` - Webhook для обработки обновлений от Telegram

## Разработка

### Структура проекта

```
telegram-sync/
├── main.go              # Точка входа приложения
├── config/
│   └── config.go        # Управление конфигурацией
├── handlers/
│   └── telegram.go      # Обработчики HTTP запросов
├── services/
│   ├── telegram.go      # Логика работы с Telegram API
│   └── github.go        # Логика работы с GitHub API
├── .env.example         # Пример файла конфигурации
├── .gitignore          # Игнорируемые файлы
└── README.md           # Документация
```

## Безопасность

- Никогда не коммитьте файл `.env` с реальными токенами
- Используйте HTTPS для webhook URL
- Храните GitHub token в безопасном месте
- Рекомендуется использовать GitHub token с ограниченными правами (только на нужный репозиторий)

## Лицензия

MIT

## Автор

hoterpoter
