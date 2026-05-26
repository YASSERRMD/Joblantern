package voice_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yasserrmd/joblantern/internal/voice"
)

func TestWhisper_Transcribe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			http.Error(w, "want multipart", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"text":"Is this Dubai job legit?"}`))
	}))
	defer srv.Close()
	a := voice.NewWhisperASR(srv.URL)
	text, err := a.Transcribe(context.Background(), []byte("RIFF....WAVE"), "en")
	if err != nil || text != "Is this Dubai job legit?" {
		t.Fatalf("text=%q err=%v", text, err)
	}
}

func TestPiper_Synthesize(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		_, _ = w.Write([]byte("RIFFdata"))
	}))
	defer srv.Close()
	p := voice.NewPiperTTS(srv.URL)
	audio, mime, err := p.Synthesize(context.Background(), "hello", "en_US-amy-low")
	if err != nil || string(audio) != "RIFFdata" || mime != "audio/wav" {
		t.Fatalf("audio=%q mime=%q err=%v", audio, mime, err)
	}
}
