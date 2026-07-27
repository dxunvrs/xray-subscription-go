package main

import (
	"log"
	"xray-subscription-go/internal/config"
)

func main() {
	config := config.Load()

	log.Println("Переменные env загружены")
}
