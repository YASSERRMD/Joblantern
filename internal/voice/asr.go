// Package voice wires Joblantern to local speech-to-text and
// text-to-speech servers. We deliberately keep llama.cpp / whisper.cpp
// / piper out of the Go binary — they are large native deps that
// break the single-binary distroless deploy model. Instead, we POST to
// HTTP endpoints (whisper.cpp server, piper-http) the operator runs
// as separate sidecar containers.
//
// Default ports follow the upstream projects:
//
//   - Whisper.cpp server  http://localhost:8090/inference
//   - Piper HTTP          http://localhost:5000/api/tts
//
// Operators wire these in deploy/voice/docker-compose.yml.
package voice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

// ASR is a speech-to-text backend.
type ASR interface {
	Transcribe(ctx context.Context, audioWAV []byte, language string) (string, error)
}

// TTS is a text-to-speech backend.
type TTS interface {
	Synthesize(ctx context.Context, text, voice string) ([]byte, string, error) // returns bytes + mime
}

// WhisperASR talks to a whisper.cpp `server` endpoint.
type WhisperASR struct {
	URL        string
	HTTPClient *http.Client
}

// NewWhisperASR — URL like "http://localhost:8090/inference".
func NewWhisperASR(url string) *WhisperASR {
	return &WhisperASR{
		URL:        url,
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// Transcribe POSTs the audio (WAV/PCM) as multipart `file`.
func (w *WhisperASR) Transcribe(ctx context.Context, audio []byte, language string) (string, error) {
	if len(audio) == 0 {
		return "", errors.New("voice: empty audio")
	}
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="file"; filename="clip.wav"`)
	hdr.Set("Content-Type", "audio/wav")
	part, err := mw.CreatePart(hdr)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(audio); err != nil {
		return "", err
	}
	if language != "" {
		_ = mw.WriteField("language", language)
	}
	_ = mw.WriteField("response_format", "json")
	_ = mw.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("whisper: %d %s", resp.StatusCode, string(raw))
	}
	var out struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("whisper decode: %w", err)
	}
	return out.Text, nil
}

// PiperTTS posts to a Piper HTTP server's /api/tts.
type PiperTTS struct {
	URL        string
	HTTPClient *http.Client
}

// NewPiperTTS — URL like "http://localhost:5000/api/tts".
func NewPiperTTS(url string) *PiperTTS {
	return &PiperTTS{
		URL:        url,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Synthesize POSTs JSON {text, voice} and returns wav bytes.
func (p *PiperTTS) Synthesize(ctx context.Context, text, voice string) ([]byte, string, error) {
	if text == "" {
		return nil, "", errors.New("voice: empty text")
	}
	body, _ := json.Marshal(map[string]string{"text": text, "voice": voice})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, "", fmt.Errorf("piper: %d %s", resp.StatusCode, string(raw))
	}
	audio, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20)) // 50MB cap
	if err != nil {
		return nil, "", err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "audio/wav"
	}
	return audio, mime, nil
}
