package main

import (
	"application"
	"context"
	"fmt"
	"time"
)

func main() {
	// Initialize the application's container.
	app := application.NewContainer()
	authClient := app.InfrastructureContainer.Get().AuthClient.Get()
	config := app.Config.Get()

	defer func() {
		if err := authClient.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configuration for the request.
	issuer := config.AuthServer.Issuer
	scopes := []string{"read", "write"}

	jwtToken, err := authClient.GenerateToken(ctx, issuer, scopes)
	if err != nil {
		fmt.Printf("could not generate token: %v", err)
	}
	fmt.Println(jwtToken)
}
