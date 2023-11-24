package load

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// LoadConfiguration is for loading necessary configs
func LoadConfiguration() (*Config, error) {
	// locate
	var config Config
	absPath, err := filepath.Abs("config/config.json")
	if err != nil {
		log.Fatalf("Could not load config file from path:%s\n", absPath)
		return nil, err
	}
	configFile, err := os.Open(absPath)
	if err != nil {
		log.Fatalf("Could not open config file from path:%s\n", absPath)
		return nil, err
	}

	// load
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return &config, err
}
