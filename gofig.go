// Gofig provides methods to load configuration values from a given JSON file.
package gofig

import (
	"encoding/json"
	"io/ioutil"
)

// Config type contains methods to fetch configuration values from a JSON file.
type Config struct {
	values map[string]interface{}
}

// Load creates a new Config type based on the supplied path to the JSON file.
func Load(path string) (*Config, error) {
	if b, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		c := &Config{}
		if err = json.Unmarshal(b, &c.values); err != nil {
			return nil, err
		}
		return c, nil
	}
}
