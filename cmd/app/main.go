package main

import (
	"fmt"
	"log"
	"xray-subscription-go/internal/config"
	"xray-subscription-go/internal/service"
)

func main() {
	config := config.Load()
	log.Println("Переменные env загружены")

	vlessBuilder := service.NewVlessLinkBuilder(config)

	link := vlessBuilder.BuildVlesLink("aboba", "hui")
	fmt.Println(link)
}
