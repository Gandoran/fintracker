package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LLM struct {
		Model       string  `yaml:"model"`
		URL         string  `yaml:"url"`
		Temperature float32 `yaml:"temperature"`
		//TODO aggiungere gestione token
	} `yaml:"llm"`

	SCRAPER struct {
		IntervalMinutes int      `yaml:"interval_minutes"`
		Feeds           []string `yaml:"feeds"`
	} `yaml:"scraper"`
}

func Load(filepath string) (*Config, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the Config file: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(fileBytes, &cfg); err != nil {
		return nil, fmt.Errorf("errore nel decodificare lo YAML: %v", err)
	}
	return &cfg, nil
}
