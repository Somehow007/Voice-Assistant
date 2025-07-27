package service

import (
	"Voice-Assistant/pkg/llm"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type FunctionResult struct {
	Result interface{} `json:"result"`
}

type AddressResp struct {
	Info     string    `json:"info"`
	Geocodes []Geocode `json:"geocodes"`
}
type Geocode struct {
	CityCode string `json:"citycode"`
	Adcode   string `json:"adcode"`
	City     string `json:"city"`
}

type WeatherResp struct {
	Info  string `json:"info"`
	Lives []Live `json:"lives"`
}

type Live struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	AdCode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Winddirection    string `json:"winddirection"`
	WindPower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	ReportTime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	HumidityFloat    string `json:"humidity_float"`
}

var ToolRegistry = map[string]func(llm.FunctionCall) (string, error){
	"get_current_time":    GetCurrentTime,
	"get_current_weather": GetCurrentWeather,
}

func CallToolByName(name string, args llm.FunctionCall) (string, error) {
	fn, ok := ToolRegistry[name]
	if !ok {
		return "", fmt.Errorf("未注册的工具: %s", name)
	}
	return fn(args)
}

func GetCurrentTime(_ llm.FunctionCall) (string, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func GetCurrentWeather(args llm.FunctionCall) (string, error) {
	_ = godotenv.Load(".env.development")
	apiKey := os.Getenv("WEATHER_API")
	addressURL := os.Getenv("WEATHER_ADDRESS_URL")
	weatherURL := os.Getenv("WEATHER_URL")

	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args.Arguments), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return "", fmt.Errorf("缺少 location 参数")
	}
	getAddressUrl := addressURL + "?key=" + apiKey + "&address=" + location
	resp, err := http.Get(getAddressUrl)
	if err != nil {
		return "", fmt.Errorf("请求地理编码API失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取地理编码响应体失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误状态码 %d : %s", resp.StatusCode, string(body))
	}
	var addressResp AddressResp
	if err := json.Unmarshal(body, &addressResp); err != nil {
		return "", fmt.Errorf("解析地理编码响应JSON失败: %w", err)
	}
	geocodes := addressResp.Geocodes
	var result strings.Builder
	for _, g := range geocodes {
		adCode := g.Adcode
		if adCode == "" {
			return "", fmt.Errorf("城市编码为空")
		}
		getWeatherUrl := weatherURL + "?key=" + apiKey + "&city=" + adCode
		resp2, err := http.Get(getWeatherUrl)
		if err != nil {
			return "", fmt.Errorf("请求天气API失败: %w", err)
		}
		defer resp2.Body.Close()
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			return "", fmt.Errorf("读取天气响应体失败: %w", err)
		}
		var weatherResp WeatherResp
		if err := json.Unmarshal(body2, &weatherResp); err != nil {
			return "", fmt.Errorf("解析天气响应JSON失败: %w", err)
		}
		if len(weatherResp.Lives) > 0 {
			live := weatherResp.Lives[0]
			result.WriteString(fmt.Sprintf("%s今天天气%s, 温度%s;", location, live.Weather, live.TemperatureFloat))
		}
	}
	return result.String(), nil
}
