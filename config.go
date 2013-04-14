package gofig

import (
	"encoding/json"
	"log"
	"os"
)

//Private vars
const (
	FILE_NAME string = "gofig.json"
)

//Types
type config struct {
	values map[string]string
}

//Public functions
func (c *config) Populate() {
	configFile, err := os.Open(FILE_NAME)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		configFile.Close()
	}()

	configFileStats, err := configFile.Stat()
	if err != nil {
		log.Panic(err)
	}

	buffer := make([]byte, configFileStats.Size())
	_, err = configFile.Read(buffer)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(buffer, &c.values)
	if err != nil {
		log.Panic(err)
	}
}

func (c *config) Val(id string) interface{} {
	return c.values[id]
}
