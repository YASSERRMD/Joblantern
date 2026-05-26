package main

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/yasserrmd/joblantern/internal/botcore"
)

// runTelegram is the long-poll Telegram adapter.
//
// Conversation model is intentionally simple — one chat, one in-flight
// listing at a time. Users can:
//
//	/start         show help
//	/verify <txt>  one-shot: submit text and wait for the verdict
//	/set ...       set submission fields (company, jurisdiction)
//	/go            submit whatever has been accumulated
//	/status <id>   re-fetch a previous verdict
//	/forget        wipe the local session
func runTelegram(token string, api *botcore.APIClient, sessions *botcore.Sessions, rl *botcore.RateLimiter, viewBase string) error {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("connect telegram: %w", err)
	}
	slog.Info("telegram adapter started", "username", bot.Self.UserName)

	updates := bot.GetUpdatesChan(tg.UpdateConfig{Offset: 0, Timeout: 30})
	for u := range updates {
		go handleUpdate(bot, api, sessions, rl, viewBase, u)
	}
	return nil
}

func handleUpdate(bot *tg.BotAPI, api *botcore.APIClient, sessions *botcore.Sessions, rl *botcore.RateLimiter, viewBase string, u tg.Update) {
	if u.Message == nil || u.Message.Chat == nil {
		return
	}
	chatID := strconv.FormatInt(u.Message.Chat.ID, 10)
	if !rl.Allow(chatID) {
		reply(bot, u.Message.Chat.ID, "Rate-limited, please slow down.")
		return
	}

	text := strings.TrimSpace(u.Message.Text)
	st := sessions.Get(chatID)
	st.UpdatedAt = time.Now()

	switch {
	case strings.HasPrefix(text, "/start"), strings.HasPrefix(text, "/help"):
		reply(bot, u.Message.Chat.ID, helpText())
	case strings.HasPrefix(text, "/verify"):
		body := strings.TrimSpace(strings.TrimPrefix(text, "/verify"))
		if body == "" {
			reply(bot, u.Message.Chat.ID, "Send: /verify <paste the recruiter message>")
			return
		}
		st.Submission.ListingText = body
		submitAndReply(bot, api, st, u.Message.Chat.ID, viewBase)
	case strings.HasPrefix(text, "/set "):
		applySet(st, strings.TrimPrefix(text, "/set "), bot, u.Message.Chat.ID)
	case text == "/go":
		if st.Submission.ListingText == "" && st.Submission.CompanyName == "" {
			reply(bot, u.Message.Chat.ID, "Nothing to submit. Use /verify or /set first.")
			return
		}
		submitAndReply(bot, api, st, u.Message.Chat.ID, viewBase)
	case strings.HasPrefix(text, "/status"):
		id := strings.TrimSpace(strings.TrimPrefix(text, "/status"))
		if id == "" {
			id = st.LastID
		}
		if id == "" {
			reply(bot, u.Message.Chat.ID, "No verification id provided and none on file.")
			return
		}
		rec, err := api.Get(context.Background(), id)
		if err != nil {
			reply(bot, u.Message.Chat.ID, "Error: "+err.Error())
			return
		}
		reply(bot, u.Message.Chat.ID, botcore.FormatVerdict(rec, viewBase+"/verifications/"+id))
	case text == "/forget":
		sessions.Reset(chatID)
		reply(bot, u.Message.Chat.ID, "Session wiped.")
	default:
		// Treat unstructured text as a listing.
		st.Submission.ListingText = text
		submitAndReply(bot, api, st, u.Message.Chat.ID, viewBase)
	}
}

func submitAndReply(bot *tg.BotAPI, api *botcore.APIClient, st *botcore.State, chatID int64, viewBase string) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	reply(bot, chatID, "Verifying…")
	id, err := api.Verify(ctx, st.Submission)
	if err != nil {
		reply(bot, chatID, "Error: "+err.Error())
		return
	}
	st.LastID = id
	rec, err := api.Wait(ctx, id, 60*time.Second)
	if err != nil {
		reply(bot, chatID, "Started "+id+" but waiting timed out. /status "+id+" to retry.")
		return
	}
	reply(bot, chatID, botcore.FormatVerdict(rec, viewBase+"/verifications/"+id))
}

func applySet(st *botcore.State, body string, bot *tg.BotAPI, chatID int64) {
	parts := strings.SplitN(body, "=", 2)
	if len(parts) != 2 {
		reply(bot, chatID, "Usage: /set field=value (company, jurisdiction, role, domain, email, phone)")
		return
	}
	key := strings.TrimSpace(strings.ToLower(parts[0]))
	val := strings.TrimSpace(parts[1])
	switch key {
	case "company":
		st.Submission.CompanyName = val
	case "jurisdiction", "country":
		st.Submission.Jurisdiction = strings.ToUpper(val)
	case "role":
		st.Submission.Role = val
	case "domain":
		st.Submission.Domain = val
	case "email":
		st.Submission.RecruiterEmail = val
	case "phone":
		st.Submission.RecruiterPhone = val
	default:
		reply(bot, chatID, "Unknown field: "+key)
		return
	}
	reply(bot, chatID, "Set "+key+" = "+val)
}

func reply(bot *tg.BotAPI, chatID int64, text string) {
	msg := tg.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	_, _ = bot.Send(msg)
}

func helpText() string {
	return strings.Join([]string{
		"Joblantern — verify job listings.",
		"",
		"Quick start: just paste a recruiter message.",
		"",
		"Commands:",
		"  /verify <text>     verify pasted text",
		"  /set field=value   accumulate fields (company, jurisdiction, role, domain, email, phone)",
		"  /go                submit accumulated fields",
		"  /status [id]       re-fetch a verdict",
		"  /forget            clear my session",
		"",
		"Joblantern is not a lawyer. Always verify with the regulator.",
	}, "\n")
}
