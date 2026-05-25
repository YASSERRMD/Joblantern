package agent

import (
	"context"
	"sync"

	"github.com/yasserrmd/joblantern/internal/mcpclient"
)

// ClientFn is a deferred mcpclient constructor. The orchestrator owns
// the lifecycle of every MCP client; sub-agents only borrow them.
type ClientFn func(ctx context.Context, server string) (*mcpclient.Client, error)

// MCPSubagent is a generic sub-agent that drives one or more MCP
// servers and converts their structured-content into Facts. Concrete
// sub-agents (address, registry, pattern, salary_law, routing) are
// wired in the binary's main.go and share this struct.
type MCPSubagent struct {
	NameStr string
	Run_    func(ctx context.Context, sub Submission) []Fact
}

func (m *MCPSubagent) Name() string                                   { return m.NameStr }
func (m *MCPSubagent) Run(ctx context.Context, sub Submission) []Fact { return m.Run_(ctx, sub) }

// runParallel is a convenience helper sub-agents use to call multiple
// MCP tools concurrently and append their results.
func runParallel(funcs ...func() []Fact) []Fact {
	var (
		mu  sync.Mutex
		out []Fact
		wg  sync.WaitGroup
	)
	for _, fn := range funcs {
		wg.Add(1)
		go func(f func() []Fact) {
			defer wg.Done()
			facts := f()
			mu.Lock()
			out = append(out, facts...)
			mu.Unlock()
		}(fn)
	}
	wg.Wait()
	return out
}
