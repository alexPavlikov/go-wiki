package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Timeout  time.Duration `mapstructure:"timeout"`
	LogLevel string        `mapstructure:"loglevel"`
	Server   Server        `mapstructure:"server"`
}

type Server struct {
	Path string `mapstructure:"path"`
	Port int    `mapstructure:"port"`
}

func (s *Server) ToString() string {
	return fmt.Sprintf("%s:%d", s.Path, s.Port)
}

func Load() (*Config, error) {
	path, fileName := fetchConfigPath()

	if path == "" || fileName == "" {
		return nil, errors.New("path or filename is empty")
	}

	var cfg Config

	cfg, err := initViper(path, fileName, cfg)
	if err != nil {
		return nil, err
	}

	level := cfg.GetLogLevel()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	slog.SetDefault(log)

	return &cfg, nil
}

func initViper(filename, path string, cfg Config) (Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config err: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config err: %w", err)
	}

	return cfg, nil
}

// Функция для определения какой файл конфига читать (local, dev, prod)
func fetchConfigPath() (path string, file string) {

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.StringVar(&file, "config_file", "", "config file name")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if file == "" {
		file = os.Getenv("CONFIG_FILE")
	}

	return path, file
}

func (cfg *Config) GetLogLevel() slog.Level {
	switch cfg.LogLevel {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
