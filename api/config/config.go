package config

import (
	"os"
)

// AppAddress app host without a protocol
var AppAddress string

// DbUser app host without a protocol
var DbUser string

// DbPsw address of
var DbPsw string

// DbHost JSON schemas dir
var DbHost string

// SetEnvs read in envs
func SetEnvs() {

	AppAddress = os.Getenv("APP_ADDRESS")
	if AppAddress == "" {
		AppAddress = "127.0.0.1:8011"
	}

	DbUser = os.Getenv("DB_USER")
	if DbUser == "" {
		DbUser = "api_user"
	}

	DbPsw = os.Getenv("DB_PSW")
	if DbPsw == "" {
		DbPsw = "123456"
	}

	DbHost = os.Getenv("DB_HOST")
	if DbPsw == "" {
		DbPsw = "0.0.0.0"
	}
}
