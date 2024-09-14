# KBOT Telegram Weather Bot

This README is available in multiple languages:
- [English](README.md) - curent page
- [Українська](README.ua.md)

KBOT is a multifunctional Telegram bot developed in Go. It provides weather information, current time in selected cities, and integrates with OpenTelemetry for metrics collection.

## Usage

1. Find the bot in Telegram: [@kbot_seo_bot](https://t.me/kbot_seo_bot).
2. Start the bot by sending the `/start` command.
3. Use the following commands to interact:
   - `Hello` — Greet the bot.
   - `Help` — Get the list of available commands.
   - `Time` — Get the current server time.
   - `Weather` — Get weather information (the bot will prompt you to enter a city name).
   - City names (Kyiv, New York, London, Seattle, Sydney) — Get the current time and weather for these cities.

## Installation

### Requirements:
- Go 1.16+
- Telegram bot token (as environment variable `TELE_TOKEN`)
- OpenWeatherMap API key (as environment variable `openweather_api_key`)
- (Optional) OpenTelemetry metrics endpoint (as environment variable `METRICS_HOST`)

### Installation Instructions:

1. Clone the repository:
   ```sh
   git clone https://github.com/YuriiKosiy/kbot.git
   cd kbot
   ```

2. Set up the required environment variables:
   ```sh
   export TELE_TOKEN="YOUR_TELEGRAM_TOKEN"
   export openweather_api_key="YOUR_OPENWEATHER_API_KEY"
   export METRICS_HOST="YOUR_METRICS_HOST" # Optional
   ```

3. Run the bot:
   ```sh
   go run main.go start
   ```

## Build & Docker

### Build (via Makefile):
```sh
make build
./kbot start
```

### Docker Build:
1. Build the image:
   ```sh
   docker build -t kbot .
   ```

2. Get the `<image-id>` by running:
   ```sh
   docker images
   ```
   The `<image-id>` is listed in the `IMAGE ID` column of the output.

3. Run the container:
   ```sh
   docker run -e TELE_TOKEN="YOUR_TELEGRAM_TOKEN" -e openweather_api_key="YOUR_OPENWEATHER_API_KEY" <image-id>
   ```

## Monitoring & Metrics

The bot supports metrics collection via OpenTelemetry. If you wish to enable this, ensure that the `METRICS_HOST` environment variable is set, and the bot will automatically connect to the remote telemetry server for metrics collection.

## How to Get a Telegram Bot Token

1. In Telegram, search for 'BotFather'.
2. Create a new bot and obtain the token.
3. Use this token as the `TELE_TOKEN` environment variable.

## How to Get an OpenWeatherMap API Key

1. Register on OpenWeatherMap: [https://openweathermap.org](https://openweathermap.org).
2. Create a new API key in your account.
3. Use this key as the `openweather_api_key` environment variable.