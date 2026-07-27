package main

import (
	"fmt"
	"log"
	"xray-subscription-go/internal/config"
	"xray-subscription-go/internal/database"
	"xray-subscription-go/internal/repository"
	"xray-subscription-go/internal/service"
)

func main() {
	config := config.Load()
	log.Println("Переменные env загружены")

	db := database.InitDB("data/app.db")
	log.Println("БД SQLite подключена")

	userRepo := repository.NewUserRepository(db)
	vlessBuilder := service.NewVlessLinkBuilder(config)

	link := vlessBuilder.BuildVlesLink("aboba", "hui")
	fmt.Println(link)
	fmt.Println(userRepo)
}
