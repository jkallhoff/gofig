/*
Gofig provides methods to load configuration values from a given JSON file.

When inflating structs, key name underscores are dropped & the first letter of the next
word is capitalized.  For example, a json key of 'my_config_value' will be
populated into a struct field called MyConfigValue.
*/
package gofig

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

// Config type contains methods to fetch configuration values from a JSON file.
type Config struct {
	values map[string]interface{}
	keyRx  *regexp.Regexp
}

// Load creates a new Config type based on the supplied path to the JSON file.
func Load(path string) (*Config, error) {
	if b, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		var c *Config
		if c, err = initConfig(b); err != nil {
			return nil, err
		}
		return c, nil
	}
}

func initConfig(rawJson []byte) (c *Config, err error) {
	rx, err := regexp.Compile(`^.{1}|\_.{1}`)
	if err == nil {
		c = &Config{keyRx: rx}
		err = json.Unmarshal(rawJson, &c.values)
	}
	return
}
