package config

import (

	"os"

	"gopkg.in/yaml.v3"
)

type Camera struct {
	Name   string
	Source interface{}
}
type Config struct {
	Cameras []Camera
}
type Event struct {
	Camera     string
	Time       float64
	Event      string
	Confidence float64
}

func CreateDefaultConfig() error {
	defaultConfig := `cameras:
  - name: cam1
    source: 0
`

	return os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
}
func CheckConfigFile() bool {
	_, err := os.Stat("config.yaml")
	if err != nil {
		return false
	}
	return true
}
func ReadConfig() (*Config, error) {
	cameras := &Config{}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cameras)
	if err != nil {
		return nil, err
	}
	return cameras, nil
}
func LoadConfig() (*Config, error) {

	if !CheckConfigFile() {
		err := CreateDefaultConfig()
		if err != nil {
			return nil, err
		}
	}
	cameras, err := ReadConfig()
	if err != nil {

		return nil, err
	}

	return cameras, nil
}
