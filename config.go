package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	strava "github.com/strava/go.strava"
)

type BoundingBox struct {
	Min, Max strava.Location
}

func (b BoundingBox) Within(position strava.Location) bool {
	return position[0] >= b.Min[0] &&
		position[0] <= b.Max[0] &&
		position[1] >= b.Min[1] &&
		position[1] <= b.Max[1]
}

type Config struct {
	GearID    string                 `json:"gear_id,omitempty"`
	Locations map[string]BoundingBox `json:"locations"`
}

func (c Config) GetLocation(position strava.Location) string {
	for name, location := range c.Locations {
		if location.Within(position) {
			return name
		}
	}

	return ""
}

func LoadConfig(path string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	if len(config.Locations) < 1 {
		return Config{}, errors.New("config contains no locations")
	}

	return config, nil
}
