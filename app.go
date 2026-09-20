package main

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/0xh4ty/quailfs/internal/keys"
	"sync"
)

type App struct {
	ctx context.Context

	mu sync.RWMutex

	unlocked bool

	userID            []byte
	ed25519PrivateKey []byte
	ed25519PublicKey  []byte
	x25519PrivateKey  []byte
	x25519PublicKey   []byte
}

type UserInfo struct {
	UserID    string
	PublicKey string
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

func (a *App) GenerateRecoveryPhrase() (string, error) {
	return keys.GenerateMnemonic(), nil
}

func (a *App) Unlock(recoveryPhrase string) (bool, error) {
	mnemonicSeed := keys.DeriveSeed(recoveryPhrase)

	ed25519Seed, err := keys.DeriveEd25519Seed(mnemonicSeed)
	if err != nil {
		return false, err
	}

	ed25519PrivateKey, ed25519PublicKey := keys.DeriveEd25519Keypair(ed25519Seed)

	x25519Seed, err := keys.DeriveX25519Seed(mnemonicSeed)
	if err != nil {
		return false, err
	}

	x25519PrivateKey, x25519PublicKey, err := keys.DeriveX25519Keypair(x25519Seed)
	if err != nil {
		return false, err
	}

	userID := keys.DeriveUserID(ed25519PublicKey)

	a.mu.Lock()
	defer a.mu.Unlock()

	a.userID = userID
	a.ed25519PrivateKey = ed25519PrivateKey
	a.ed25519PublicKey = ed25519PublicKey
	a.x25519PrivateKey = x25519PrivateKey
	a.x25519PublicKey = x25519PublicKey
	a.unlocked = true

	return true, nil
}

func (a *App) GetUserInfo() (UserInfo, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.unlocked {
		return UserInfo{}, errors.New("identity is locked")
	}

	return UserInfo{
		UserID:    hex.EncodeToString(a.userID),
		PublicKey: hex.EncodeToString(a.ed25519PublicKey),
	}, nil
}
