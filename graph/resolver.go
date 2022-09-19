package graph

import (
	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Config struct {
	*resolver.Services
}

type Resolver struct {
	*resolver.Services
}

func NewResolver(cfg Config) (*Resolver, error) {
	return &Resolver{Services: cfg.Services}, nil
}
