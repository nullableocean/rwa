package services

import (
	"rwa/internal/models"
	"rwa/pkg/passwordcryptor"
	"time"
)

type UserRepository interface {
	GetById(id int64) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Save(models.User) (int64, error)
	Update(models.User) error
	Delete(models.User) error
}

type UserService struct {
	userRepo  UserRepository
	passCrypt passwordcryptor.PasswordCryptor
}

func NewUserService(userRepo UserRepository, passCryptor passwordcryptor.PasswordCryptor) *UserService {
	return &UserService{
		userRepo:  userRepo,
		passCrypt: passCryptor,
	}
}

func (us *UserService) CreateUser(info models.UserCreateInfo) (*models.User, error) {
	if err := info.Validate(); err != nil {
		return nil, err
	}

	hashedPass, err := us.getPasswordHash(info.Password)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	user := models.User{
		Email:          info.Email,
		Username:       info.Username,
		Bio:            info.Bio,
		Image:          info.Image,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		HashedPassword: hashedPass,
	}

	id, err := us.userRepo.Save(user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return &user, nil
}

func (us *UserService) UpdateUser(user models.User, newInfo models.UserUpdateInfo) (*models.User, error) {
	err := newInfo.Validate()
	if err != nil {
		return nil, err
	}

	updated := time.Now()

	user.UpdatedAt = updated
	user.Email = newInfo.Email
	user.Bio = newInfo.Bio
	user.Username = newInfo.Username
	user.Image = newInfo.Image

	err = us.userRepo.Update(user)
	return &user, err
}

func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := us.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserById(id int64) (*models.User, error) {
	user, err := us.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) VerificatePassword(user models.User, password string) bool {
	return us.passCrypt.CheckHash(password, user.HashedPassword)
}

func (us *UserService) getPasswordHash(password string) (string, error) {
	return us.passCrypt.Crypt(password)
}
