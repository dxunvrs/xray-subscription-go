package service

import (
	"fmt"
	"log"

	"xray-subscription-go/internal/repository"
)

type XraySyncService struct {
	userRepo       *repository.UserRepository
	xrayGrpcClient *XrayGrpcClient
}

func NewXraySyncService(userRepo *repository.UserRepository, xrayGrpcClient *XrayGrpcClient) *XraySyncService {
	return &XraySyncService{
		userRepo:       userRepo,
		xrayGrpcClient: xrayGrpcClient,
	}
}

func (s *XraySyncService) SyncUsersOnStartup() error {
	log.Println("Начинаем синхронизацию пользователей с XRay...")

	users, err := s.userRepo.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch users from db for sync: %w", err)
	}

	successCount := 0
	failCount := 0

	for _, user := range users {
		err := s.xrayGrpcClient.AddUser(user.Email, user.UUID)
		if err != nil {
			log.Printf("Ошибка добавления пользователя %s (%s): %v", user.Email, user.UUID, err)
			failCount++
			continue
		}

		successCount++
	}

	log.Printf("Синхронизация завершена. Успешно: %d, Ошибок: %d", successCount, failCount)
	return nil
}
