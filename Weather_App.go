package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const API_KEY = " " 

type WeatherResponse struct {
	Name string `json:"name"` //City name
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp     float64 `json:"temp"`    
		Humidity int     `json:"humidity"` 
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`      
		Description string `json:"description"` 
	} `json:"weather"`
	Cod int `json:"cod"`
}


func getWeather(cityName string) (*WeatherResponse, error) {
	// Construct the API URL with city and API key
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cityName, API_KEY)


	httpResponse, requestError := http.Get(url)
	if requestError != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", requestError)
	}
	defer httpResponse.Body.Close()

	var weather WeatherResponse
	if requestError := json.NewDecoder(httpResponse.Body).Decode(&weather); requestError != nil {
		return nil, fmt.Errorf("error decoding response: %v", requestError)
	}

	// Checking the API returned a valid city
	if weather.Cod != 200 {
		return nil, fmt.Errorf("city '%s' not found or API error", cityName)
	}

	return &weather, nil
}

func getWeatherIcon(condition string) string {
	switch strings.ToLower(condition) {
	case "clear":
		return "☀️"
	case "clouds":
		return "☁️"
	case "rain":
		return "🌧️"
	case "drizzle":
		return "🌦️"
	case "thunderstorm":
		return "⛈️"
	case "snow":
		return "❄️"
	case "mist", "fog", "haze":
		return "🌫️"
	default:
		return "🌍"
	}
}

// Main function of the program
func main() {
	// greeting for the user
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("🌤️  Hello! Welcome to the Weather App!")

	for {
		fmt.Print("\nEnter city or cities: ") //Ask the user to input
		userInput, _ := inputReader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		formattedInput := strings.ReplaceAll(userInput, ",", " ") // Use can use comma or space to separate inputs
		cities := strings.Fields(formattedInput)

		for _, cityName := range cities {
			cityName = strings.TrimSpace(cityName)
			if cityName == "" { //Skip if there empty string
				continue
			}

			fmt.Printf("\nFetching weather for %s...\n", cityName)
			weatherData, requestError := getWeather(cityName)
			if requestError != nil {
				fmt.Printf("❌ Error: %v\n", requestError) // If the fetching weather is fail disply error message
				continue
			}

			//Display the weather details
			icon := getWeatherIcon(weatherData.Weather[0].Main)
			fmt.Printf("\n🌆 %s, %s\n", weatherData.Name, weatherData.Sys.Country)
			fmt.Printf("Temperature: %.1f°C\n", weatherData.Main.Temp)
			fmt.Printf("Humidity: %d%%\n", weatherData.Main.Humidity)
			fmt.Printf("Conditions: %s %s\n", weatherData.Weather[0].Description, icon)
			fmt.Println(strings.Repeat("-", 30))
		}

		
		fmt.Print("\nDo you want to check another city? (yes/no): ")
		choice, _ := inputReader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice != "yes" && choice != "y" {
			fmt.Println("👋 Thank you for using the Go Weather App. Stay safe!")
			break
		}
	}
}








