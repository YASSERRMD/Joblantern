// Package config loads runtime configuration for the Joblantern
// process. Concerns are split: mcp.go knows about MCP-server wiring;
// future files will cover LLM-provider wiring, observability, and so
// on.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// MCPTransport is the on-the-wire transport the agent uses to reach
// an MCP server.
type MCPTransport string

const (
	TransportStdio MCPTransport = "stdio"
	TransportHTTP  MCPTransport = "http"
)

// MCPServer is one entry in config/mcp.yaml.
type MCPServer struct {
	Name        string            `yaml:"name"`
	Enabled     bool              `yaml:"enabled"`
	Transport   MCPTransport      `yaml:"transport"`
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args"`
	URL         string            `yaml:"url"`
	Env         map[string]string `yaml:"env"`
	CallTimeout time.Duration     `yaml:"call_timeout"`
	MaxRetries  int               `yaml:"max_retries"`
}

// MCPConfig is the parsed config/mcp.yaml document.
type MCPConfig struct {
	Servers []MCPServer `yaml:"servers"`
}

// LoadMCP reads and validates the MCP registry from disk.
func LoadMCP(path string) (*MCPConfig, error) {
	data, err := os.ReadFile(path) //nolint:gosec // operator-supplied path
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var cfg MCPConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return &cfg, nil
}

func (c *MCPConfig) validate() error {
	seen := make(map[string]struct{}, len(c.Servers))
	for i, s := range c.Servers {
		if s.Name == "" {
			return fmt.Errorf("servers[%d]: name is required", i)
		}
		if _, dup := seen[s.Name]; dup {
			return fmt.Errorf("duplicate server name: %s", s.Name)
		}
		seen[s.Name] = struct{}{}

		switch s.Transport {
		case TransportStdio:
			if s.Command == "" {
				return fmt.Errorf("servers[%d] %q: stdio transport requires command", i, s.Name)
			}
		case TransportHTTP:
			if s.URL == "" {
				return fmt.Errorf("servers[%d] %q: http transport requires url", i, s.Name)
			}
		case "":
			return fmt.Errorf("servers[%d] %q: transport is required", i, s.Name)
		default:
			return fmt.Errorf("servers[%d] %q: unknown transport %q", i, s.Name, s.Transport)
		}
	}
	return nil
}

// EnabledByName returns the enabled server with the given name, or nil
// if it is not configured or is disabled.
func (c *MCPConfig) EnabledByName(name string) *MCPServer {
	for i := range c.Servers {
		s := &c.Servers[i]
		if s.Name == name && s.Enabled {
			return s
		}
	}
	return nil
}
