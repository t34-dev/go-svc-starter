package main

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/logger"
	"github.com/t34-dev/go-svc-starter/internal/service"
	auth_service "github.com/t34-dev/go-svc-starter/internal/service/auth"
	role_manager "github.com/t34-dev/go-svc-starter/pkg/role-manager"
	"github.com/t34-dev/go-utils/pkg/etcd"
	"github.com/t34-dev/go-utils/pkg/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	role_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/role"
	session_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/session"
	user_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/transaction"
)

type AuthService struct {
	service service.Service
}

func NewService(pool *pgxpool.Pool) service.Service {
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

	serv := service.Service{}
	deps := service.NewDeps(serv, repos, nil, role_manager.NewRoleManager(nil), txManager)
	serv.Auth = auth_service.New(deps, []byte("KEY"))

	return serv
}

func main() {
	ctx := context.Background()
	logs.Init(logger.GetCore(zap.NewAtomicLevelAt(zap.InfoLevel), "logs/XXX.log"))

	// ETCD
	cli, err := etcd.NewClient(clientv3.Config{
		Endpoints: []string{"localhost:2378"},
	}, nil)
	if err != nil {
		logs.Fatal("failed to create etcd client", zap.Error(err))
	}
	accessManager := role_manager.NewRoleManager(cli)
	err = accessManager.UpdateConfigsFromEtcd(ctx)
	if err != nil {
		logs.Fatal("failed to update ETCD config", zap.Error(err))
	}
	// Update configuration only if it has changed in the storage
	err = accessManager.WatchConfig(ctx, func(err2 error, key string, newValue []byte) {
		if err2 != nil {
			logs.Error("failed to update watch config", zap.Error(err2))
			return
		}
		logs.Warn("updated watch config", zap.String("key", key), zap.String("newValue", string(newValue)))

		// user
		ok, err := accessManager.CheckAccess("user", "blog", "sex")
		fmt.Println("user [blog] sex", ok, err)
		// user
		ok, err = accessManager.CheckAccess("user", "blog", "write")
		fmt.Println("user [blog] write", ok, err)
		// user
		ok, err = accessManager.CheckAccess("user", "blog", "read")
		fmt.Println("user [blog] read", ok, err)
	})

	fmt.Println("==== ROLES =====")

	// Database connection string
	connString := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	// Create a connection pool
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Check the connection
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	srv := NewService(pool)

	// Usage example
	email := "test@example.com"
	username := "testuser"
	password := "testpassword"
	deviceKey := "device_key"
	userAgent := "Test User Agent"

	// User registration
	res, err := srv.Auth.Registration(ctx, email, username, password, deviceKey, userAgent)
	if err != nil {
		log.Printf("Failed Registration: %v", err)
	}
	fmt.Printf("1) Registration successful. Tokens: %+v\n", res)

	// User login
	res, err = srv.Auth.Login(ctx, email, password, deviceKey, userAgent)
	if err != nil {
		log.Fatalf("Failed Login: %v", err)
	}
	fmt.Printf("2) Login successful. Tokens: %+v\n", res)

	// Token validation
	validateRes, err := srv.Auth.ValidateToken(ctx, res.Token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	fmt.Printf("3) Token validation result: %+v\n", validateRes)

	// Get user information
	userID, _ := uuid.Parse(validateRes.UserID)
	userInfo, err := srv.Auth.GetUserInfo(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("4) User Info: %+v\n", userInfo)

	// Get all roles
	roles, err := srv.Auth.GetAllRoles(ctx)
	if err != nil {
		log.Fatalf("Failed to get all roles: %v", err)
	}
	fmt.Printf("5) All roles: %+v\n", roles)

	// Create a new role
	newRoleID, err := srv.Auth.CreateRole(ctx, "NewRole")
	if err != nil {
		log.Fatalf("Failed to create new role: %v", err)
	}
	fmt.Printf("6) Created new role with ID: %d\n", newRoleID)

	// Assign role to user
	err = srv.Auth.AssignRoleToUser(ctx, userID, newRoleID)
	if err != nil {
		log.Fatalf("Failed to assign role to user: %v", err)
	}
	fmt.Printf("7) Assigned role %d to user %s\n", newRoleID, userID)

	// Token refresh
	newRes, err := srv.Auth.RefreshToken(ctx, res.RefreshToken)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	fmt.Printf("8) Refreshed tokens: %+v\n", newRes)

	// Get active user sessions
	sessions, err := srv.Auth.GetActiveSessions(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get active sessions: %v", err)
	}
	fmt.Printf("9) Active sessions: %+v\n", sessions)

	// Revoke session
	//if len(sessions) > 0 {
	//	err = srv.Auth.RevokeSession(ctx, sessions[0].ID)
	//	if err != nil {
	//		log.Fatalf("Failed to revoke session: %v", err)
	//	}
	//	fmt.Printf("10) Revoked session %s\n", sessions[0].ID)
	//}
	//
	// User logout
	//err = srv.Auth.Logout(ctx, newRes.Token)
	//if err != nil {
	//	log.Fatalf("Failed to logout: %v", err)
	//}
	//fmt.Println("11) Logout successful")
	//
	// Delete role
	//err = srv.Auth.DeleteRole(ctx, newRoleID)
	//if err != nil {
	//	log.Fatalf("Failed to delete role: %v", err)
	//}
	//fmt.Printf("12) Deleted role %d\n", newRoleID)

	// Start cleanup of inactive sessions in the background
	go func() {
		for {
			time.Sleep(24 * time.Hour) // Run every 24 hours
			_ = srv.Auth.CleanupInactiveSessions(ctx)
		}
	}()

	// Here you can add code to start an HTTP server or other application logic

	// Infinite loop to keep the program running
	select {}
}
