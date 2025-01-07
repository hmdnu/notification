package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	SiakadUrl string
	SlcUrl    string
	Nim       string
	Password  string
}

var Env env

func init() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .Env file")
	}

	Env.SiakadUrl = os.Getenv("SIAKAD_URL")
	Env.Nim = os.Getenv("NIM")
	Env.Password = os.Getenv("PASSWORD")
	Env.SlcUrl = os.Getenv("SIAKAD_SLC")
}
