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
	VaultProvider        string `env:"VAULT_PROVIDER" env-default:"hashicorp"` // azure/hashicorp/aws
	VaultURL             string `env:"VAULT_URL"`
	VaultUsername        string `env:"VAULT_USERNAME"` // only for hashicorp provider
	VaultPassword        string `env:"VAULT_PASSWORD"` // only for hashicorp provider
	HashicorpTransitPath string `env:"HASHICORP_TRANSIT_PATH" env-default:"transit"`
	AWSAccessKeyID       string `env:"AWS_ACCESS_KEY_ID"` // only for aws provider
	AWSSecretKey         string `env:"AWS_SECRET_KEY"`    // only for aws provider
	DBURI                string `env:"DB_URI" env-required`
	Port                 int    `env:"FP_PORT" env-default:"7003"`
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
