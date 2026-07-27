package service

import (
	"xray-subscription-go/internal/dto"
	"xray-subscription-go/internal/model"
	"xray-subscription-go/internal/repository"
	"xray-subscription-go/internal/util"

	"github.com/google/uuid"
)

type XrayGrpcService interface {
	AddUser(email, userUUID string) error
	RemoveUser(email string) error
	GetUserUplink(email string) int64
	GetUserDownlink(email string) int64
}

type UserService struct {
	userRepo   *repository.UserRepository
	xrayClient XrayGrpcService
}

func NewUserService(userRepo *repository.UserRepository, xrayClient XrayGrpcService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		xrayClient: xrayClient,
	}
}

func (s *UserService) CreateUser(email string) (*dto.UserResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	userUUID := uuid.New().String()

	user := &model.User{
		Email: email,
		UUID:  userUUID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if err := s.xrayClient.AddUser(email, userUUID); err != nil {
		return nil, err
	}

	return s.toUserDto(user), nil
}

func (s *UserService) DeleteUser(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.DeleteByEmail(email); err != nil {
		return err
	}

	return s.xrayClient.RemoveUser(email)
}

func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, *s.toUserDto(&user))
	}

	return responses, nil
}

func (s *UserService) FindUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.toUserDto(user), nil
}

func (s *UserService) toUserDto(user *model.User) *dto.UserResponse {
	uplink := s.xrayClient.GetUserUplink(user.Email)
	downlink := s.xrayClient.GetUserDownlink(user.Email)
	total := uplink + downlink

	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		UUID:  user.UUID,
		Traffic: dto.UserTraffic{
			UplinkFormatted:   util.FormatBytes(uplink),
			DownlinkFormatted: util.FormatBytes(downlink),
			TotalFormatted:    util.FormatBytes(total),
			UplinkBytes:       uplink,
			DownlinkBytes:     downlink,
			TotalBytes:        total,
		},
	}
}
