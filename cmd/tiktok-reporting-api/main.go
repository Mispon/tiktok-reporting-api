package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mispon/tiktok-reporting-api/internal/env"
	"github.com/mispon/tiktok-reporting-api/internal/handlers"
)

func main() {
	environment := env.New()

	ctx := context.Background()

	fmt.Printf("starting with env: %v\n", environment)
	if err := run(ctx, &environment); err != nil {
		log.Fatal(err)
	}
}

func run(context context.Context, env *env.Env) error {
	handler := handlers.New(env)
	if err := handler.Init(context); err != nil {
		return err
	}

	fmt.Printf("TikTok reporting API listening on %s\n", env.Endpoint)
	return http.ListenAndServe(env.Endpoint, nil)
}
