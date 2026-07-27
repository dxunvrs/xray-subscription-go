package service

import "xray-subscription-go/internal/repository"

type SubscriptionService struct {
	userRepo     *repository.UserRepository
	vlessBuilder *VlessLinkBuilder
}

func NewXraySubscriptionService(userRepo *repository.UserRepository, vlessBuilder *VlessLinkBuilder) *SubscriptionService {
	return &SubscriptionService{
		userRepo:     userRepo,
		vlessBuilder: vlessBuilder,
	}
}

func (s *SubscriptionService) GetSubscription(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	vlessLink := s.vlessBuilder.BuildVlesLink(user.UUID, user.Email)
	return vlessLink, nil
}
