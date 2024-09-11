package auth_imp

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	descAuth "github.com/t34-dev/go-svc-starter/pkg/api/auth_v1"
)

func (s *ImplementedAuth) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {

	fmt.Println(color.RedString("%20-s", "Note info:"), color.GreenString("Привет как дела?"))
	return &descAuth.LoginResponse{}, nil
}
