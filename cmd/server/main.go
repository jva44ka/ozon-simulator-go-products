//go:generate swag init -g cmd/server/main.go --dir ./internal,./cmd
package main

import (
	"fmt"
	"os"

	app2 "github.com/jva44ka/ozon-simulator-go-products/internal/app"
)

const (
	grpcPort = ":50051"
	httpPort = ":8081"
)

func main() {
	fmt.Println("app starting")

	app, err := app2.NewApp(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}
