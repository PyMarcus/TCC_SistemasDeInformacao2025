package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct{
	DatabaseUrl string 
}

// LoadConfig: loads database and api settings from .env file
// including the database connection URL (DatabaseUrl).
func LoadConfig(dotenvFilePath string) (*Config, error){

	if dotenvFilePath == ""{
		dotenvFilePath = "../../../.env"
	}
	
	err := godotenv.Load(dotenvFilePath)

	if err != nil{
		log.Println("[-] Fail to load .env " + err.Error())
		return nil, err 
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == ""{
		log.Println("Database url is required")
		return nil, errors.New("DatabaseUrl is missing")
	}

	return &Config{
		DatabaseUrl: dbUrl,
	}, nil
}