package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	AuthToken string `json:"auth_token"`
}

func LoadOrCreateConfig() (*Config, error) {
	c := &Config{}
	configFile := configFile()
	err := loadFromFile(configFile, c)
	if err != nil {
		c.AuthToken = promptForToken()
		err := c.Save()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func configFile() string {
	configFile := os.Getenv("PBL_CONFIG")
	if configFile != "" {
		return configFile
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal("failed to determin the current user.")
	}
	return filepath.Join(usr.HomeDir, ".config", "pbl")
}

func promptForToken() (pass string) {
	pass = os.Getenv("PINBOARD_AUTH_TOKEN")
	if pass != "" {
		return
	}

	fmt.Printf("Pinboard auth token: ")
	fmt.Scanln(&pass)

	return
}

func (c *Config) Save() error {
	return saveToFile(configFile(), c)
}

func loadFromFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	for {
		if err := decoder.Decode(v); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

func saveToFile(filename string, v interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(v)
}
