package main

import (
	"context"
	"sync"
)

type App struct {
	ctx context.Context

	mu sync.RWMutex

	unlocked bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ctx = ctx
}

func (a *App) IsUnlocked() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.unlocked
}
