package main

import (
	"flag"
	"github.com/k3env/grafana-webhook/config"
	"github.com/k3env/grafana-webhook/templates"
	"github.com/k3env/grafana-webhook/tg"
	"github.com/k3env/grafana-webhook/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

var (
	path = flag.String("path", "hook.yaml", "Path to config file")
)

func main() {
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "02.01.06 15:04:05"})
	cfg, err := config.LoadConfig(*path)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	log.Info().Str("path", cfg.TemplatesDirectory).Msg("Loading templates")
	tpls := templates.Load(cfg.TemplatesDirectory)

	sender := tg.NewTelegramSender(cfg.Telegram, tpls)

	srv := http.NewServeMux()
	srv.Handle("POST /webhook", web.NewWebhookHandler(sender))
	srv.Handle("GET /templates", web.NewTemplatesHandler(tpls))

	log.Info().Str("addr", cfg.ListenAddress).Msg("Starting server")
	err = http.ListenAndServe(cfg.ListenAddress, srv)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}
