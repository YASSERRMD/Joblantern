package agent

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Store is the (in-memory in v1) backing for verification requests +
// results. A future Postgres-backed implementation can satisfy the same
// interface; sqlc-generated queries (Phase 02) are ready for it.
type Store interface {
	Create(ctx context.Context, sub Submission) (string, error)
	Get(ctx context.Context, id string) (*Record, error)
	List(ctx context.Context, limit int) ([]*Record, error)
	Save(ctx context.Context, rec *Record) error
}

// Record is the persisted state of one verification.
type Record struct {
	ID         string     `json:"id"`
	Submission Submission `json:"submission"`
	Status     string     `json:"status"` // pending | running | completed | failed
	Verdict    *Verdict   `json:"verdict,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MemoryStore is a non-persistent Store useful for tests and Phase 13's
// initial agent endpoints. Phase 13 also adds a Postgres-backed Store
// behind the same interface.
type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]*Record
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[string]*Record)}
}

// Create generates an id (if absent), stores the submission, returns the id.
func (s *MemoryStore) Create(_ context.Context, sub Submission) (string, error) {
	id := sub.ID
	if id == "" {
		id = uuid.NewString()
	}
	now := time.Now().UTC()
	rec := &Record{
		ID:         id,
		Submission: sub,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	rec.Submission.ID = id
	s.mu.Lock()
	s.records[id] = rec
	s.mu.Unlock()
	return id, nil
}

func (s *MemoryStore) Get(_ context.Context, id string) (*Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := s.records[id]
	if r == nil {
		return nil, nil
	}
	cp := *r
	return &cp, nil
}

func (s *MemoryStore) List(_ context.Context, limit int) ([]*Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Record, 0, len(s.records))
	for _, r := range s.records {
		cp := *r
		out = append(out, &cp)
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (s *MemoryStore) Save(_ context.Context, rec *Record) error {
	rec.UpdatedAt = time.Now().UTC()
	cp := *rec
	s.mu.Lock()
	s.records[rec.ID] = &cp
	s.mu.Unlock()
	return nil
}
