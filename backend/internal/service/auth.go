package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidInput       = errors.New("invalid input")
)

const demoEmail = "auditor@secaudit.local"
const demoPassword = "AuditorDemo!2026"
const demoName = "Demo Auditor"

type AuthService struct {
	users *repository.UserRepo
}

func NewAuthService(users *repository.UserRepo) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if len(name) < 2 || len(name) > 80 || len(password) < 8 || len(password) > 72 {
		return nil, ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil || strings.Contains(email, " ") {
		return nil, ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{Email: email, Name: name, PasswordHash: string(hash), Role: models.RoleAuditor}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.users.ByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*models.User, error) {
	return s.users.ByID(ctx, userID)
}

func (s *AuthService) Demo(ctx context.Context) (*models.User, error) {
	user, err := s.users.ByEmail(ctx, demoEmail)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	return s.Register(ctx, demoName, demoEmail, demoPassword)
}
