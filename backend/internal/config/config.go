package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	LLM      LLMConfig
	AsrTts   AsrTtsConfig
	Amap     AmapConfig
}

type ServerConfig struct {
	Port string
	Host string
	Mode string // gin模式: debug, release, tes
}

type DatabaseConfig struct {
	DBPath string // SQLite数据库文件路径
}

type LLMConfig struct {
	APIKey               string
	TextEndpoint         string
	FunctionCallEndpoint string
	Model                string
}

type AsrTtsConfig struct {
	AppId       string
	SecretId    string
	SecretKey   string
	AsrEndpoint string
	TtsEndpoint string
	Speaker     string
}

type AmapConfig struct {
	AddressUrl string
	WeatherUrl string
	ApiKey     string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			DBPath: getEnv("DB_PATH", "data.db"),
		},
		LLM: LLMConfig{
			APIKey:               getEnv("LLM_API_KEY", ""),
			TextEndpoint:         getEnv("LLM_TEXT_ENDPOINT", ""),
			FunctionCallEndpoint: getEnv("LLM_FUNCTION_CALLING_ENDPOINT", ""),
			Model:                getEnv("LLM_MODEL", "qwen-plus"),
		},
		AsrTts: AsrTtsConfig{
			AppId:       getEnv("APPID", ""),
			SecretId:    getEnv("ASR_TTS_SECRET_ID", ""),
			SecretKey:   getEnv("ASR_TTS_SECRET_KEY", ""),
			AsrEndpoint: getEnv("ASR_ENDPOINT", ""),
			TtsEndpoint: getEnv("TTS_ENDPOINT", ""),
			Speaker:     getEnv("TTS_SPEAKER", ""),
		},
		Amap: AmapConfig{
			AddressUrl: getEnv("AMAP_ADDRESS_URL", ""),
			WeatherUrl: getEnv("AMAP_WEATHER_URL", ""),
			ApiKey:     getEnv("AMAP_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err != nil {
			return intValue
		}
	}
	return defaultValue
}
