package main

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/config"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	err, resultChan := config.New(ctx, "configs/dev.yaml", ".env")
	err, resultChan = config.New(ctx, "configs/dev.yaml", ".env")
	err, resultChan = config.New(ctx, "configs/dev.yaml", ".env")
	if err != nil {
		panic(err)
	}
	fmt.Println("DATA:", config.App().Name(), config.Grpc().Port())

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case result := <-resultChan:
				if result.Error != nil {
					fmt.Println("Ошибка из Watch:", result.Error)
				} else {
					fmt.Println("Успешно обновили")
				}
			case <-ticker.C:
				fmt.Println("FROM:", config.App().Name(), config.Grpc().Port())
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("END")
		ticker.Stop()
	}
}
