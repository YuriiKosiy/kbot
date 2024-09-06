package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	telebot "gopkg.in/telebot.v3"
)

// UserSession tracks the user's context, such as if they are awaiting city input.
type UserSession struct {
	AwaitingCity bool
}

var (
	// Maps user sessions
	userSessions = struct {
		sync.RWMutex
		sessions map[int64]*UserSession
	}{sessions: make(map[int64]*UserSession)}

	// Environment variables
	TelegramToken     = os.Getenv("TELE_TOKEN")
	OpenWeatherAPIKey = os.Getenv("openweather_api_key")
	MetricsHost       = os.Getenv("METRICS_HOST") // Assuming this is defined
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {
	if MetricsHost == "" {
		// Логуємо попередження, але не припиняємо виконання програми
		log.Println("Warning: METRICS_HOST не встановлено; пропуск ініціалізації метрик")
		return
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		fmt.Printf("Failed to create exporter: %v\n", err)
		panic(err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second)),
		),
	)

	otel.SetMeterProvider(mp)
}

func pmetrics(ctx context.Context, payload string) {
	meter := otel.GetMeterProvider().Meter("kbot_commands_counter")

	counter, err := meter.Int64Counter(fmt.Sprintf("kbot_command_%s", payload))
	if err != nil {
		fmt.Printf("Error creating counter: %v\n", err)
		return
	}
	counter.Add(ctx, 1)
}

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Telegram bot that provides various functionalities.",
	Long:    "A versatile Telegram bot built with Go, providing weather information, current times in selected cities, and more.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerodriver.NewProductionLogger()

		fmt.Printf("kbot %s started", appVersion)

		if TelegramToken == "" || OpenWeatherAPIKey == "" {
			log.Fatalf("Missing Telegram token or OpenWeather API key.")
		}

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TelegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			logger.Fatal().Str("Error", err.Error()).Msg("Please check TELE_TOKEN")
			return
		} else {
			logger.Info().Str("Version", appVersion).Msg("kbot started")
		}

		menu := &telebot.ReplyMarkup{
			ReplyKeyboard: [][]telebot.ReplyButton{
				{{Text: "Hello"}, {Text: "Help"}},
				{{Text: "Time"}, {Text: "Weather"}},
				{{Text: "Kyiv"}, {Text: "Boston"}, {Text: "London"}},
				{{Text: "Vienna"}, {Text: "Tbilisi"}, {Text: "Vancouver"}},
			},
		}

		kbot.Handle("/start", func(m telebot.Context) error {
			logger.Info().Str("Payload", m.Text()).Msg(m.Message().Payload)
			return m.Send("Welcome to Kbot! Use the menu for available commands.", menu)
		})

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			logger.Info().Str("Payload", m.Text()).Msg(m.Message().Payload)

			payload := m.Text()
			userID := m.Sender().ID

			pmetrics(context.Background(), payload)

			userSessions.RLock()
			session, exists := userSessions.sessions[userID]
			userSessions.RUnlock()

			if exists && session.AwaitingCity {
				city := payload
				if city == "" {
					return m.Send("Please provide a valid city name.")
				}
				weatherInfo, err := getWeatherInfo(OpenWeatherAPIKey, city)
				if err != nil {
					return m.Send(fmt.Sprintf("Failed to get weather information: %s", err))
				}
				userSessions.Lock()
				session.AwaitingCity = false
				userSessions.Unlock()
				return m.Send(weatherInfo)
			}

			switch payload {
			case "Hello":
				err := m.Send(fmt.Sprintf("Hi! I'm Kbot %s! I can help you with time and weather.", appVersion))
				return err
			case "Help":
				err := m.Send("Available commands: Time, Weather, or select a city for the current time.")
				return err
			case "Time":
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				err := m.Send(fmt.Sprintf("Current server time is: %s", currentTime))
				return err
			case "Weather":
				userSessions.Lock()
				userSessions.sessions[userID] = &UserSession{AwaitingCity: true}
				userSessions.Unlock()
				return m.Send("Please enter the city name.")
			case "Kyiv", "New_York", "London", "Seattle", "Sydney":
				err := m.Send(fmt.Sprintf("Current time in %s: %s", payload, getTime(payload)))
				return err
			default:
				err := m.Send("Unknown command. Use the menu for available commands.")
				return err
			}
		})

		kbot.Start()
	},
}

// getWeatherInfo requests weather data from OpenWeather API.
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

// getTime returns the current time in the specified location.
func getTime(location string) string {
	var locName string
	switch location {
	case "Kyiv":
		locName = "Europe/Kiev"
	case "New_York":
		locName = "America/New_York"
	case "London":
		locName = "Europe/London"
	case "Seattle":
		locName = "America/Seattle"
	case "Sydney":
		locName = "Australia/Sydney"
	default:
		return "Invalid location"
	}

	loc, err := time.LoadLocation(locName)
	if err != nil {
		return "Invalid location"
	}
	return time.Now().In(loc).Format("15:04:05")
}

func init() {
	ctx := context.Background()
	initMetrics(ctx)

	rootCmd.AddCommand(kbotCmd)
}
