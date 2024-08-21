# KBOT Telegram Weather Bot

This Telegram bot provides weather information for a specified city. It uses the Telegram API and OpenWeatherMap API to fetch the latest weather data.

## Usage

1. Open Telegram and search for the bot [@kbot_seo_bot](https://t.me/kbot_seo_bot).
2. Start interacting with the bot by sending the `/start` command.
3. Use the following commands to interact with the bot:
   - `/start` - Start the bot and get a welcome message.
   - `/help` - Show a list of available commands.
   - `/hello` - Greet the bot.
   - `/time` - Get the current server time.
   - `/weather` - Get weather information (after sending this command, enter the name of the city).

## Installation

To install and run this bot on your own system, follow these steps:

### Requirements

- Go 1.16 or later
- Telegram bot token
- OpenWeatherMap API key

### Steps

1. Clone the repository:
   ```sh
   git clone https://github.com/YuriiKosiy/kbot.git
   cd kbot

2. Create a config.json file based on the config.json.sample:
    ```sh
    cp config.json.sample config.json

3. Edit the config.json file to include your tokens:
     ```sh
    "TELE_TOKEN": "YOUR_TELE_TOKEN",
    "openweather_api_key": "YOUR_OPENWEATHER_API_KEY"

### How to Get a Telegram Bot Token

1. Open Telegram and search for BotFather.
2. Create a new bot by following BotFather's instructions.
3. Obtain your bot token from BotFather and paste it into the config.json file instead "YOUR_TELE_TOKEN".

### How to Get an OpenWeatherMap API Key

1. Sign up at OpenWeatherMap > https://openweathermap.org/.
2. Go to your profile and create a new API key.
3. Paste your API key into the config.json file instead "YOUR_OPENWEATHERMAP_API_KEY".

## Run your bot
   ```sh
   ./kbot start
   ```

## Build kbot (from Makefile)
   ```sh
   make build
   ./kbot start
   ```

## Docker build (Dockerfile)
   ```sh
   docker build .
   docker run sha256:<image sha256 from last step docker build>
   ```
