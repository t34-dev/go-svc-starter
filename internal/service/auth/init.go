package auth_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("")

var _ service.AuthService = &authService{}

type authService struct {
	opt service.Options
}

func New(opt service.Options, secretKey []byte) service.AuthService {
	jwtKey = secretKey
	return &authService{
		opt: opt,
	}
}

func (s *authService) Registration(ctx context.Context, email, username, password, deviceKey, deviceName string) (*model.AuthTokens, error) {
	// Проверяем, существует ли пользователь
	_, err := s.opt.Repos.User.GetUserByLogin(ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	_, err = s.opt.Repos.User.GetUserByLogin(ctx, username)
	if err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var authTokens *model.AuthTokens

	err = s.opt.TxManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Вставляем нового пользователя
		userID, errTx := s.opt.Repos.User.CreateUser(ctx, email, username, string(hashedPassword))
		if errTx != nil {
			return errTx
		}

		// Создаем сессию
		sessionID, errTx := s.opt.Repos.Session.CreateSession(ctx, userID, deviceKey, deviceName)
		if errTx != nil {
			return errTx
		}

		// Генерируем токены
		token, refreshToken, errTx := generateTokens(userID, sessionID)
		if errTx != nil {
			return errTx
		}

		authTokens = &model.AuthTokens{
			Token:        token,
			RefreshToken: refreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return authTokens, nil
}

func (s *authService) Login(ctx context.Context, login, password, deviceKey, deviceName string) (*model.AuthTokens, error) {
	// Находим пользователя
	user, err := s.opt.Repos.User.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Создаем или обновляем сессию
	sessionID, err := s.opt.Repos.Session.UpsertSession(ctx, user.ID, deviceKey, deviceName)
	if err != nil {
		return nil, err
	}

	// Генерируем новые токены
	token, refreshToken, err := generateTokens(user.ID, sessionID)
	if err != nil {
		return nil, err
	}

	return &model.AuthTokens{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	sessionID, err := uuid.Parse(claims["session_id"].(string))
	if err != nil {
		return err
	}

	return s.opt.Repos.Session.DeleteSession(ctx, sessionID)
}

func (s *authService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	user, err := s.opt.Repos.User.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	sessions, err := s.opt.Repos.Session.GetActiveSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %v", err)
	}

	roles, err := s.opt.Repos.Role.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	user.Roles = roles

	return &model.UserInfo{
		User:     user,
		Sessions: sessions,
	}, nil
}

func (s *authService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]model.Session, error) {
	return s.opt.Repos.Session.GetActiveSessions(ctx, userID)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokens, error) {
	// Парсим refresh token
	claims, err := validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}
	sessionID, err := uuid.Parse(claims["session_id"].(string))
	if err != nil {
		return nil, err
	}

	// Проверяем существование сессии в базе данных
	session, err := s.opt.Repos.Session.GetSessionByID(ctx, sessionID)
	if err != nil || session.UserID != userID {
		return nil, errors.New("invalid session")
	}

	// Обновляем время последнего использования сессии
	err = s.opt.Repos.Session.UpdateSessionLastUsed(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Генерируем новые токены
	newAccessToken, newRefreshToken, err := generateTokens(userID, sessionID)
	if err != nil {
		return nil, err
	}

	return &model.AuthTokens{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*model.ValidateTokenResponse, error) {
	claims, err := validateToken(token)
	if err != nil {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	sessionID, err := uuid.Parse(claims["session_id"].(string))
	if err != nil {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	// Проверяем существование сессии в базе данных
	session, err := s.opt.Repos.Session.GetSessionByID(ctx, sessionID)
	if err != nil || session.UserID != userID {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	return &model.ValidateTokenResponse{
		Valid:  true,
		UserID: userID.String(),
	}, nil
}

func (s *authService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.opt.Repos.Session.DeleteSession(ctx, sessionID)
}

func (s *authService) CleanupInactiveSessions(ctx context.Context) error {
	err := s.opt.Repos.Session.CleanupInactiveSessions(ctx)
	return err
}
