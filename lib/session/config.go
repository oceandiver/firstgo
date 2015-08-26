package session

import (
       "encoding/json"
       "io/ioutil"
)

type LetsgoConfig struct {
     RootDir      string             `json:"rootDir"`
}

func LoadConfiguration(jsonFile string) (*LetsgoConfig, error) {
     content, e := ioutil.ReadFile(jsonFile)

     if e != nil {
     	return nil, e
	}
	return loadConfigFromString(content)
}

func loadConfigFromString(content []byte) (*LetsgoConfig, error) {

     configuration := &LetsgoConfig{}
     e := json.Unmarshal(content, configuration)

     return configuration, e
}
