package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mispon/tiktok-reporting-api/internal/job"

	"github.com/mispon/tiktok-reporting-api/internal/env"
	"github.com/mispon/tiktok-reporting-api/internal/handlers"
)

func main() {
	appEnv := env.New()
	loadAdvertsFromFile(&appEnv)

	ctx := context.WithValue(context.Background(), "env", &appEnv)

	runJob(ctx, &appEnv)

	fmt.Printf("starting with env: %v\n", appEnv)
	if err := runApi(ctx, &appEnv); err != nil {
		log.Fatal(err)
	}
}

func runJob(ctx context.Context, env *env.Env) {
	scheduler, err := job.New(ctx, env)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := scheduler.Schedule(); err != nil {
			log.Fatal(err)
		}
	}()
}

func runApi(context context.Context, env *env.Env) error {
	handler := handlers.New()
	if err := handler.Init(context, env); err != nil {
		return err
	}

	fmt.Printf("TikTok reporting API listening on %s\n", env.Endpoint)
	return http.ListenAndServe(env.Endpoint, nil)
}

func loadAdvertsFromFile(env *env.Env) {
	file, err := os.Open("advert_ids.txt")
	if err != nil {
		fmt.Println("advert_ids.txt not found")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var advertIds []int64
	for scanner.Scan() {
		if id, err := strconv.ParseInt(scanner.Text(), 10, 64); err == nil {
			advertIds = append(advertIds, id)
		}
	}

	env.AdvertiserIds = advertIds
}
