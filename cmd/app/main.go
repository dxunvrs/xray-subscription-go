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

type MockXrayService struct{}

func (m *MockXrayService) AddUser(email, userUUID string) error { return nil }
func (m *MockXrayService) RemoveUser(email string) error        { return nil }
func (m *MockXrayService) GetUserUplink(email string) int64     { return 1024 * 1024 * 50 } // 50 MB
func (m *MockXrayService) GetUserDownlink(email string) int64   { return 1024 * 1024 }

func main() {
	config := config.Load()
	log.Println("Переменные env загружены")

	db := database.InitDB("data/app.db")
	log.Println("БД SQLite подключена")

	userRepo := repository.NewUserRepository(db)
	vlessBuilder := service.NewVlessLinkBuilder(config)

	xrayMock := &MockXrayService{}
	userService := service.NewUserService(userRepo, xrayMock)
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
