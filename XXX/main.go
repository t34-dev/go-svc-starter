package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	role_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/role"
	session_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/session"
	user_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/transaction"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	txManager db.TxManager
	repos     repository.Repository
}

var jwtKey = []byte("your-secret-key") // В реальном приложении используйте безопасный метод хранения ключа

func NewAuthService(pool *pgxpool.Pool) *AuthService {
	dbClient, err := pg.New(pool, nil)
	if err != nil {
		log.Fatalln(err)
	}
	txManager := transaction.NewTransactionManager(dbClient.DB())
	repos := repository.Repository{
		User:    user_repository.New(dbClient),
		Session: session_repository.New(dbClient),
		Role:    role_repository.New(dbClient),
	}
	return &AuthService{
		txManager: txManager,
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

		// Создаем сессию
		sessionID, errTx := s.repos.Session.CreateSession(ctx, userID, deviceKey, deviceName)
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

	// Создаем или обновляем сессию
	sessionID, err := s.repos.Session.UpsertSession(ctx, user.ID, deviceKey, deviceName)
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

func (s *AuthService) Logout(ctx context.Context, token string) error {
	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	sessionID, err := uuid.Parse(claims["session_id"].(string))
	if err != nil {
		return err
	}

	return s.repos.Session.DeleteSession(ctx, sessionID)
}

func (s *AuthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	user, err := s.repos.User.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	sessions, err := s.repos.Session.GetActiveSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %v", err)
	}

	roles, err := s.repos.Role.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	user.Roles = roles

	return &model.UserInfo{
		User:     user,
		Sessions: sessions,
	}, nil
}

func (s *AuthService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]model.Session, error) {
	return s.repos.Session.GetActiveSessions(ctx, userID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokens, error) {
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
	session, err := s.repos.Session.GetSessionByID(ctx, sessionID)
	if err != nil || session.UserID != userID {
		return nil, errors.New("invalid session")
	}

	// Обновляем время последнего использования сессии
	err = s.repos.Session.UpdateSessionLastUsed(ctx, sessionID)
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

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*model.ValidateTokenResponse, error) {
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
	session, err := s.repos.Session.GetSessionByID(ctx, sessionID)
	if err != nil || session.UserID != userID {
		return &model.ValidateTokenResponse{Valid: false}, nil
	}

	return &model.ValidateTokenResponse{
		Valid:  true,
		UserID: userID.String(),
	}, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.repos.Session.DeleteSession(ctx, sessionID)
}

func (s *AuthService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.repos.Role.GetAllRoles(ctx)
}

func (s *AuthService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	return s.repos.Role.AssignRoleToUser(ctx, userID, roleID)
}

func (s *AuthService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	return s.repos.Role.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *AuthService) CreateRole(ctx context.Context, roleName string) (int64, error) {
	return s.repos.Role.CreateRole(ctx, roleName)
}

func (s *AuthService) DeleteRole(ctx context.Context, roleID int64) error {
	return s.repos.Role.DeleteRole(ctx, roleID)
}

func (s *AuthService) UpdateRole(ctx context.Context, roleID int64, newRoleName string) error {
	return s.repos.Role.UpdateRole(ctx, roleID, newRoleName)
}

func (s *AuthService) cleanupInactiveSessions(ctx context.Context) {
	err := s.repos.Session.CleanupInactiveSessions(ctx)
	if err != nil {
		log.Printf("Error cleaning up inactive sessions: %v", err)
	}
}

func generateTokens(userID uuid.UUID, sessionID uuid.UUID) (string, string, error) {
	// Access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID.String(),
		"session_id": sessionID.String(),
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID.String(),
		"session_id": sessionID.String(),
		"exp":        time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 дней
	})

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

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
	email := "test@example.com"
	username := "testuser"
	password := "testpassword"
	deviceKey := "device_key"
	userAgent := "Test User Agent"

	// Регистрация пользователя
	res, err := authService.Registration(ctx, email, username, password, deviceKey, userAgent)
	if err != nil {
		log.Printf("Failed Registration: %v", err)
	}
	fmt.Printf("1) Registration successful. Tokens: %+v\n", res)

	// Логин пользователя
	res, err = authService.Login(ctx, email, password, deviceKey, userAgent)
	if err != nil {
		log.Fatalf("Failed Login: %v", err)
	}
	fmt.Printf("2) Login successful. Tokens: %+v\n", res)

	// Валидация токена
	validateRes, err := authService.ValidateToken(ctx, res.Token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	fmt.Printf("3) Token validation result: %+v\n", validateRes)

	// Получение информации о пользователе
	userID, _ := uuid.Parse(validateRes.UserID)
	userInfo, err := authService.GetUserInfo(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("4) User Info: %+v\n", userInfo)

	// Получение всех ролей
	roles, err := authService.GetAllRoles(ctx)
	if err != nil {
		log.Fatalf("Failed to get all roles: %v", err)
	}
	fmt.Printf("5) All roles: %+v\n", roles)

	// Создание новой роли
	newRoleID, err := authService.CreateRole(ctx, "NewRole")
	if err != nil {
		log.Fatalf("Failed to create new role: %v", err)
	}
	fmt.Printf("6) Created new role with ID: %d\n", newRoleID)

	// Назначение роли пользователю
	err = authService.AssignRoleToUser(ctx, userID, newRoleID)
	if err != nil {
		log.Fatalf("Failed to assign role to user: %v", err)
	}
	fmt.Printf("7) Assigned role %d to user %s\n", newRoleID, userID)

	// Обновление токена
	newRes, err := authService.RefreshToken(ctx, res.RefreshToken)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	fmt.Printf("8) Refreshed tokens: %+v\n", newRes)

	// Получение активных сессий пользователя
	sessions, err := authService.GetActiveSessions(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get active sessions: %v", err)
	}
	fmt.Printf("9) Active sessions: %+v\n", sessions)

	// Отзыв сессии
	//if len(sessions) > 0 {
	//	err = authService.RevokeSession(ctx, sessions[0].ID)
	//	if err != nil {
	//		log.Fatalf("Failed to revoke session: %v", err)
	//	}
	//	fmt.Printf("10) Revoked session %s\n", sessions[0].ID)
	//}
	//
	//// Выход пользователя
	//err = authService.Logout(ctx, newRes.Token)
	//if err != nil {
	//	log.Fatalf("Failed to logout: %v", err)
	//}
	//fmt.Println("11) Logout successful")
	//
	//// Удаление роли
	//err = authService.DeleteRole(ctx, newRoleID)
	//if err != nil {
	//	log.Fatalf("Failed to delete role: %v", err)
	//}
	//fmt.Printf("12) Deleted role %d\n", newRoleID)

	// Запуск очистки неактивных сессий в фоновом режиме
	go func() {
		for {
			time.Sleep(24 * time.Hour) // Запуск каждые 24 часа
			authService.cleanupInactiveSessions(ctx)
		}
	}()

	// Здесь можно добавить код для запуска HTTP-сервера или другой логики приложения

	// Бесконечный цикл, чтобы программа не завершалась
	select {}
}
