package utilities

import (
	"encoding/json"
	"os"

	"github.com/EddieAlvarez01/administrator_courses/models"
)

//GetDatabaseConfiguration Decoded json file in struct ConfigurationDB
func GetDatabaseConfiguration() (models.ConfigurationDB, error) {
	config := models.ConfigurationDB{}
	file, err := os.Open("./configuration.json")
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
