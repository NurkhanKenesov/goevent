package db

import (
	"goevent/internal/config"
	"testing"
)

func TestConnect(t *testing.T) {
	cfg := &config.Config{}
	_, _ = Connect(cfg)
}
