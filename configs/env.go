package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv("MONGOURI")
}

func PORT() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv("PORT")
}

func SecretKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv((""))
}
