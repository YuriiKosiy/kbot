# KBOT Telegram Weather Bot

KBOT — це багатофункціональний Telegram бот, написаний мовою Go. Він надає інформацію про погоду, час в обраних містах і має інтеграцію з OpenTelemetry для збору метрик.

## Використання

1. Знайдіть бота в Telegram: [@kbot_seo_bot](https://t.me/kbot_seo_bot).
2. Стартуйте бота командою `/start`.
3. Використовуйте такі команди для взаємодії:
   - `Hello` — Привітання з ботом.
   - `Help` — Отримайте список доступних команд.
   - `Time` — Отримайте поточний серверний час.
   - `Weather` — Отримайте інформацію про погоду (бот попросить ввести назву міста).
   - Назви міст (Kyiv, New York, London, Seattle, Sydney) — Отримайте поточний час та погоду в цих містах.

## Встановлення

### Вимоги:
- Go 1.16+
- Токен Telegram бота (як змінна оточення `TELE_TOKEN`)
- API-ключ OpenWeatherMap (як змінна оточення `openweather_api_key`)
- (Опціонально) Хост для збору метрик в OpenTelemetry (як змінна оточення `METRICS_HOST`)

### Інструкції встановлення:

1. Клонуйте репозиторій:
   ```sh
   git clone https://github.com/YuriiKosiy/kbot.git
   cd kbot
   ```

2. Встановіть необхідні змінні оточення:
   ```sh
   export TELE_TOKEN="YOUR_TELEGRAM_TOKEN"
   export openweather_api_key="YOUR_OPENWEATHER_API_KEY"
   export METRICS_HOST="YOUR_METRICS_HOST" # Опціонально
   ```

3. Запустіть бота:
   ```sh
   go run main.go start
   ```

## Build & Docker

### Build (через Makefile):
```sh
make build
./kbot start
```

### Docker-збірка:
1. Збірка образу:
   ```sh
   docker build -t kbot .
   ```

2. Get the `<image-id>` by running:
   ```sh
   docker images
   ```
   Ідентифікатор `<image-id>` вказано у стовпчику `IMAGE ID` виводу.

3. Запуск контейнера:
   ```sh
   docker run -e TELE_TOKEN="YOUR_TELEGRAM_TOKEN" -e openweather_api_key="YOUR_OPENWEATHER_API_KEY" <image-id>
   ```

## Моніторинг та Метрики

Бот підтримує збір метрик через OpenTelemetry. Якщо ви хочете включити цю функцію, встановіть змінну оточення `METRICS_HOST`, і бот автоматично підключиться до віддаленого сервера для збору метрик.

## Як отримати токен Telegram бота

1. У Telegram знайдіть 'BotFather'.
2. Створіть нового бота і отримайте токен.
3. Використайте цей токен, як змінну оточення `TELE_TOKEN`.

## Як отримати OpenWeatherMap API ключ

1. Зареєструйтесь на OpenWeatherMap: [https://openweathermap.org](https://openweathermap.org).
2. Створіть новий API ключ.
3. Використайте цей ключ, як змінну оточення `openweather_api_key`.
