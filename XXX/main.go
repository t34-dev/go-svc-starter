package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	session_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/sessions"
	user_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/transaction"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	txManager db.TxManager
	builder   sq.StatementBuilderType
	repos     repository.Repository
}

var jwtKey = []byte("your-secret-key") // В реальном приложении используйте безопасный метод хранения ключа

func NewAuthService(pool *pgxpool.Pool) *AuthService {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	dbClient, err := pg.New(pool, nil)
	if err != nil {
		log.Fatalln(err)
	}
	txManager := transaction.NewTransactionManager(dbClient.DB())
	repos := repository.Repository{
		Common:  nil,
		User:    user_repository.New(dbClient),
		Session: session_repository.New(dbClient),
	}
	return &AuthService{
		txManager: txManager,
		builder:   builder,
		repos:     repos,
	}
}

func (s *AuthService) Registration(ctx context.Context, email, username, password, deviceKey, deviceName string) (*model.AuthTokens, error) {
	// Проверяем, существует ли пользователь
	_, err := s.repos.User.GetUserByLogin(ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	_, err = s.repos.User.GetUserByLogin(ctx, username)
	if err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var authTokens *model.AuthTokens

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Вставляем нового пользователя
		userID, errTx := s.repos.User.CreateUser(ctx, email, username, string(hashedPassword))
		if errTx != nil {
			return errTx
		}

		// Генерируем токены
		token, refreshToken, errTx := generateTokens(userID)
		if errTx != nil {
			return errTx
		}

		// Вставляем информацию об устройстве (сессии)
		errTx = s.repos.Session.CreateSession(ctx, userID, deviceKey, deviceName, refreshToken)
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

func (s *AuthService) Login(ctx context.Context, login, password, deviceKey, deviceName string) (*model.AuthTokens, error) {
	// Находим пользователя
	user, err := s.repos.User.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Генерируем новые токены
	token, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Обновляем или создаем устройство (сессию)
	err = s.repos.Session.UpsertSession(ctx, user.ID, deviceKey, deviceName, refreshToken)
	if err != nil {
		return nil, err
	}

	return &model.AuthTokens{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	userID := int64(claims["user_id"].(float64))
	refreshToken := claims["refresh_token"].(string)

	return s.repos.Session.DeleteSession(ctx, userID, refreshToken)
}

func (s *AuthService) GetUserInfo(ctx context.Context, userID int64) (*model.UserInfo, error) {
	user, err := s.repos.User.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	currentSessionID, err := s.repos.Session.GetCurrentSession(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get current device: %v", err)
	}
	user.CurrentSessionID = currentSessionID

	devices, err := s.repos.Session.GetActiveSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %v", err)
	}

	return &model.UserInfo{
		User:     user,
		Sessions: devices,
	}, nil
}

func (s *AuthService) GetActiveSessions(ctx context.Context, userID int64) ([]model.Session, error) {
	return s.repos.Session.GetActiveSessions(ctx, userID)
}
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, deviceKey, deviceName string) (*model.AuthTokens, error) {
	// Находим устройство (сессию)
	device, err := s.repos.Session.GetSessionByRefreshToken(ctx, refreshToken, deviceKey)
	if err != nil {
		return nil, errors.New("invalid refresh token or device")
	}

	if device.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	// Генерируем новые токены
	token, newRefreshToken, err := generateTokens(device.UserID)
	if err != nil {
		return nil, err
	}

	// Обновляем устройство (сессию)
	err = s.repos.Session.UpdateSession(ctx, device.UserID, deviceKey, deviceName, newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &model.AuthTokens{
		Token:        token,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*model.ValidateTokenResponse, error) {
	claims, err := validateToken(token)
	if err != nil {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	userID := int64(claims["user_id"].(float64))
	refreshToken := claims["refresh_token"].(string)

	device, err := s.repos.Session.GetSessionByRefreshToken(ctx, refreshToken, "")
	if err != nil || device.UserID != userID {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	return &model.ValidateTokenResponse{
		Valid:  true,
		UserID: fmt.Sprintf("%v", userID),
	}, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, userID int64, deviceID string) error {
	return s.repos.Session.DeleteSession(ctx, userID, deviceID)
}

func (s *AuthService) cleanupInactiveSessions(ctx context.Context) {
	err := s.repos.Session.CleanupInactiveSessions(ctx)
	if err != nil {
		log.Printf("Error cleaning up inactive sessions: %v", err)
	}
}

func generateTokens(userID int64) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		return "", "", err
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshToken)

	return tokenString, refreshTokenString, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func main() {
	ctx := context.Background()
	// Строка подключения к базе данных
	connString := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	// Создаем пул соединений
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверка соединения
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	authService := NewAuthService(pool)

	// Пример использования
	email := gofakeit.Email()
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, true, true, 10)
	deviceKey := "device_key"
	userAgent := gofakeit.UserAgent()

	res, err := authService.Registration(ctx, email, login, password, deviceKey, userAgent)
	if err != nil {
		log.Fatalf("Failed Registration: %v", err)
	}
	fmt.Println(res)

	res, err = authService.Login(ctx, email, password, deviceKey, userAgent)
	if err != nil {
		log.Fatalf("Failed Login via EMAIL: %v", err)
	}
	fmt.Println(res)

	res, err = authService.Login(ctx, login, password, deviceKey, userAgent)
	if err != nil {
		log.Fatalf("Failed Login via EMAIL: %v", err)
	}
	fmt.Println(res)

	// Пример использования GetUserInfo
	claims, err := validateToken(res.Token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	userID := int64(claims["user_id"].(float64))

	userInfo, err := authService.GetUserInfo(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("User Info: %+v\n", userInfo)

	// Вывод информации о активных устройствах
	activeSessions, err := authService.GetActiveSessions(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get active devices: %v", err)
	}
	fmt.Printf("Active Sessions: %+v\n", activeSessions)

	// Start a goroutine to clean up inactive sessions periodically
	go func() {
		for {
			time.Sleep(24 * time.Hour) // Run once a day
			authService.cleanupInactiveSessions(ctx)
		}
	}()

	// Бесконечный цикл, чтобы программа не завершалась
	select {}
}
