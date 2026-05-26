# Joblantern bot

`cmd/joblantern-bot` is a thin chat-channel adapter in front of the
Joblantern HTTP API. v1 ships the Telegram adapter; WhatsApp
(`whatsmeow`) and IVR (Twilio / open-source PBX) adapters slot into the
same conversation engine.

## Run

```bash
export TELEGRAM_BOT_TOKEN=<from @BotFather>
export JOBLANTERN_API=http://localhost:8080      # joblantern web server
export JOBLANTERN_VIEW=http://localhost:8080      # public URL for verdict links
go run ./cmd/joblantern-bot
```

Or via compose:

```bash
TELEGRAM_BOT_TOKEN=... docker compose --profile bots up -d joblantern-bot
```

The `bots` compose profile keeps the bot opt-in so deployments without
a Telegram token are not blocked at startup.

## Commands

| Command | Effect |
|---|---|
| `/start`, `/help` | Show command list |
| `/verify <text>` | One-shot: paste a recruiter message |
| `/set field=value` | Accumulate fields (company, jurisdiction, role, domain, email, phone) |
| `/go` | Submit accumulated fields |
| `/status [id]` | Re-fetch a verdict by id |
| `/forget` | Wipe the local session |

Anything not starting with `/` is treated as listing text.

## Rate limits

Default: 10 messages per minute per chat. Configurable in a follow-up
PR (currently hard-coded in `cmd/joblantern-bot/main.go`).

## OCR (deferred)

The Phase 26 roadmap calls for image-OCR so users can forward
screenshots. That requires Tesseract via `gosseract` (cgo) — a
material change to the single-binary distroless model. Lands in a
follow-up branch with a separate Dockerfile.

## Privacy

The bot stores conversation state **in memory only**. Restarting the
process clears every session. No chat content is logged at INFO
level. See `docs/PRIVACY.md` for the server's stance.
