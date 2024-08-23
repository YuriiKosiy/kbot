package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

// Config
type Config struct {
	//	TelegramToken     string `json:"TELE_TOKEN"`
	OpenWeatherAPIKey string `json:"openweather_api_key"`
}

// loadConfig from config.json
func loadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// UserSession
type UserSession struct {
	AwaitingCity bool
}

var (
	// userSessions
	userSessions = struct {
		sync.RWMutex
		sessions map[int64]*UserSession
	}{sessions: make(map[int64]*UserSession)}
	// TeleToken bot
	TelegramToken = os.Getenv("TELE_TOKEN")
)

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("kbot %s started", appVersion)

		config, err := loadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %s", err)
			return
		}

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TelegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN in config file. %s", err)
			return
		}

		// hello
		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			log.Print(m.Message().Payload, m.Text())
			userID := m.Sender().ID

			userSessions.RLock()
			session, exists := userSessions.sessions[userID]
			userSessions.RUnlock()

			if exists && session.AwaitingCity {
				city := m.Text() // Get City
				if city == "" {
					return m.Send("Please provide a valid city name.")
				}
				weatherInfo, err := getWeatherInfo(config.OpenWeatherAPIKey, city)
				if err != nil {
					return m.Send(fmt.Sprintf("Failed to get weather information: %s", err))
				}
				userSessions.Lock()
				session.AwaitingCity = false
				userSessions.Unlock()
				return m.Send(weatherInfo)
			}

			switch m.Text() {
			case "/hello":
				err = m.Send(fmt.Sprintf("Hello I`m kbot %s!", appVersion))
			default:
				err = m.Send("Unknown command. Use /help to see available commands.")
			}
			return err
		})

		// help
		kbot.Handle("/help", func(m telebot.Context) error {
			helpText := `Available commands:
				/start - Start the bot
				/help - Show this help message
				/hello - Say hello
				/time - Get current server time
				/weather - Get weather information`
			return m.Send(helpText)
		})

		// time
		kbot.Handle("/time", func(m telebot.Context) error {
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			return m.Send(fmt.Sprintf("Current server time is: %s", currentTime))
		})

		// weather
		kbot.Handle("/weather", func(m telebot.Context) error {
			userID := m.Sender().ID
			userSessions.Lock()
			userSessions.sessions[userID] = &UserSession{AwaitingCity: true}
			userSessions.Unlock()
			return m.Send("Please enter the city name.")
		})

		// start
		kbot.Handle("/start", func(m telebot.Context) error {
			return m.Send("Welcome to kbot! Use /help to see available commands.")
		})

		kbot.Start()
	},
}

// OpenWeatherMap API
func getWeatherInfo(apiKey, city string) (string, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	log.Printf("Requesting weather data with URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed response: %s", body)
		return "", fmt.Errorf("failed to get weather data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var weatherData struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return "", err
	}

	weatherInfo := fmt.Sprintf("Current weather in %s: %s, %.2f°C",
		weatherData.Name, weatherData.Weather[0].Description, weatherData.Main.Temp)
	return weatherInfo, nil
}

func init() {
	rootCmd.AddCommand(kbotCmd)

}
