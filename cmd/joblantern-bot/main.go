// Package main is the joblantern-bot binary — a multi-channel chat
// adapter in front of the Joblantern HTTP API.
//
// v1 ships only the Telegram adapter; the structure leaves room for
// WhatsApp (whatsmeow) and IVR adapters to be added without changing
// the API client or conversation logic.
package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/yasserrmd/joblantern/internal/botcore"
)

func main() {
	adapter := flag.String("adapter", getenv("BOT_ADAPTER", "telegram"), "telegram (whatsapp + ivr TBD)")
	apiBase := flag.String("api", getenv("JOBLANTERN_API", "http://localhost:8080"), "Joblantern API base URL")
	apiKey := flag.String("api-key", getenv("JOBLANTERN_API_KEY", ""), "Joblantern API key")
	viewBase := flag.String("view", getenv("JOBLANTERN_VIEW", ""), "Public web base URL for verdict links (defaults to -api)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if *viewBase == "" {
		*viewBase = *apiBase
	}

	client := botcore.New(*apiBase, *apiKey)
	sessions := botcore.NewSessions()
	limiter := botcore.NewRateLimiter(10, time.Minute) // 10 msgs / minute / chat

	switch strings.ToLower(*adapter) {
	case "telegram":
		token := getenv("TELEGRAM_BOT_TOKEN", "")
		if token == "" {
			logger.Error("TELEGRAM_BOT_TOKEN not set")
			os.Exit(2)
		}
		if err := runTelegram(token, client, sessions, limiter, *viewBase); err != nil {
			logger.Error("telegram exit", "err", err)
			os.Exit(1)
		}
	default:
		logger.Error("unknown adapter", "adapter", *adapter)
		os.Exit(2)
	}
}

func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
