package main

import (
	"log"
	"net/http"
	"xray-subscription-go/internal/app"
	"xray-subscription-go/internal/config"
	"xray-subscription-go/internal/database"
	"xray-subscription-go/internal/handler"
	"xray-subscription-go/internal/repository"
	"xray-subscription-go/internal/service"
)

// @title           Xray Subscription API
// @version         2.0
// @description     API сервиса управления подписками Xray.
// @BasePath        /
// @securityDefinitions.basic BasicAuth
func main() {
	config := config.Load()
	log.Println("Переменные env загружены")

	db := database.InitDB("data/app.db")
	log.Println("БД SQLite подключена")

	userRepo := repository.NewUserRepository(db)
	vlessBuilder := service.NewVlessLinkBuilder(config)

	xrayService, err := service.NewXrayGrpcClient(config.XrayGrpcHost, config.XrayGrpcPort, config.XrayVlessInbound)
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC Xray: %v", err)
	}
	defer xrayService.Close()

	syncService := service.NewXraySyncService(userRepo, xrayService)
	if err := syncService.SyncUsersOnStartup(); err != nil {
		log.Printf("Ошибка при загрузке пользователей из БД: %v", err)
	}

	userService := service.NewUserService(userRepo, xrayService)
	subService := service.NewXraySubscriptionService(userRepo, vlessBuilder)

	subHandler := handler.NewSubHandler(subService)
	adminHandler := handler.NewAdminUserHandler(userService)
	router := app.NewRouter(adminHandler, subHandler, config)

	addr := ":12258"
	log.Printf("Сервер запущен")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
