package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg     _Configuration
	cfgOnce sync.Once
	envFile *string
)

type _Configuration struct {
	VaultURL string `env:"VAULT_URL" env-required`
	DBURI    string `env:"DB_URI" env-required`
	Port     int    `env:"FP_PORT" env-default:"7003"`
}

// ReadConfig reads the configuration file and sets the envFile variable
// If the file is not found, it will try to read the file from enviroment variable
func ReadConfig(file string) {
	log.Printf(`Reading config file: "%s"`, file)
	if _, err := os.Stat(file); err != nil {
		log.Fatalf("Config file is not exist")
	}
	cfgOnce.Do(func() {
		envFile = &file
		err := cleanenv.ReadConfig(file, &cfg)
		if err != nil {
			err := cleanenv.ReadEnv(&cfg)
			if err != nil {
				log.Fatalf("Config error %s", err.Error())
			}
			fileFlag := "nofile"
			envFile = &fileFlag
		}
	})
}

func Configuration() _Configuration {
	if envFile == nil {
		fmt.Printf(`Configuration file is not set. Using default .env`)
		file := ".env"
		ReadConfig(file)
	}
	err := cleanenv.UpdateEnv(&cfg)
	if err != nil {
		log.Fatalf("Config error %s", err.Error())
	}
	return cfg
}
