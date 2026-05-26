# Voice interface

Speech-to-text and text-to-speech for users who prefer to speak.

## Design

Joblantern itself does **not** embed whisper.cpp or piper. Those are
large native binaries that would break the single-binary distroless
deploy model (ADR-001). Instead, `internal/voice` is an HTTP client
that talks to:

- A **whisper.cpp `server`** running on `:8090` for transcription.
- A **piper-http** instance on `:5000` for synthesis.

Operators run these as side-car containers, paid for in CPU and disk
where they value the offline guarantee. The Joblantern binary itself
stays slim.

## Wiring (operator-side, sketch)

```yaml
# deploy/voice/docker-compose.yml — not shipped by default; see below.
services:
  whisper:
    image: ggerganov/whisper.cpp:latest
    command: ["./server", "--model", "/models/ggml-small.bin", "--host", "0.0.0.0", "--port", "8090"]
    volumes: [whisper-models:/models]
  piper:
    image: rhasspy/wyoming-piper:latest
    ports: ["5000:10200"]
    volumes: [piper-voices:/data]
```

(The exact compose file is operator-specific; we don't ship one
because the model files are large and the choice of voice is
locale-specific.)

## API (v0.2 plan)

A future PR exposes:

- `POST /api/v1/voice/verify` — multipart upload of an audio clip;
  transcribes → submits via `/api/v1/verify` → returns id.
- `GET /api/v1/voice/verdict/{id}.wav` — synthesised spoken verdict.

The Phase 26 Telegram bot will gain a voice-message handler that
proxies into these endpoints.

## Dialect notes

- Emirati Arabic — Whisper's `ar` model performs reasonably; verify
  per dialect with a corpus before claiming coverage.
- Philippine English — Whisper's `en` model is fine; consider the
  `large-v3` quantised builds for accent robustness.
- Indian English, Bengali, Urdu — same as above. Test, then tune.

## Privacy

Audio is never written to disk by the Go binary; it is held in memory
for the duration of the request and discarded. Whisper.cpp's `server`
likewise streams transcription without retaining audio (verify in
your operator deployment).
