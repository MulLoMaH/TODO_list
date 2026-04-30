package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"`  //чтение переменных окружения в поле конфинурациионной структуры и указываем что это значение обязательное
	Folder string `envconfig:"FOLDER" required:"true"` //чтение переменных окружения в поле конфинурациионной структуры и указываем что это значение обязательное
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return config
}
