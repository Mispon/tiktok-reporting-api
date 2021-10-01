package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jessevdk/go-flags"

	"github.com/mispon/tiktok-reporting-api/internal/handlers"
)

type Opts struct {
	Endpoint     string `short:"e" long:"endpoint" description:"API endpoint" default:"0.0.0.0:80"`
	AppId        int    `short:"i" long:"app_id" description:"TikTok app id"`
	AppSecret    string `short:"s" long:"app_secret" description:"TikTok app secret"`
	SandboxToken string `short:"t" long:"sandbox_token" description:"TikTok app sandbox token"`
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("starting with opts: %v\n", opts)
	if err := run(&opts); err != nil {
		log.Fatal(err)
	}
}

func run(opts *Opts) error {
	handler := handlers.New(opts.AppId, opts.AppSecret, opts.SandboxToken)
	handler.Init()

	fmt.Printf("TikTok reporting API listening on %s\n", opts.Endpoint)
	return http.ListenAndServe(opts.Endpoint, nil)
}
