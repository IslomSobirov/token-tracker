package tracker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	"token-tracker/internal/config"
)

// Transaction represents a token transaction
type Transaction struct {
	Type      string
	Amount    uint64
	Timestamp time.Time
}

// TokenTracker handles token operations
type TokenTracker struct {
	client    *rpc.Client
	wsClient  *ws.Client
	programID solana.PublicKey
	balances  map[string]uint64
	history   map[string][]Transaction
	mu        sync.RWMutex
}

// NewTokenTracker creates a new token tracker instance
func NewTokenTracker(cfg *config.Config) (*TokenTracker, error) {
	client := rpc.New(cfg.RPCEndpoint)
	wsClient, err := ws.Connect(context.Background(), cfg.WSEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	programID, err := solana.PublicKeyFromBase58(cfg.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("invalid program ID: %v", err)
	}

	return &TokenTracker{
		client:    client,
		wsClient:  wsClient,
		programID: programID,
		balances:  make(map[string]uint64),
		history:   make(map[string][]Transaction),
	}, nil
}

// Close closes the token tracker
func (tt *TokenTracker) Close() {
	if tt.wsClient != nil {
		tt.wsClient.Close()
	}
}

// Deposit handles token deposits
func (tt *TokenTracker) Deposit(ctx context.Context, address string, amount uint64) error {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tt.balances[address] += amount
	tt.history[address] = append(tt.history[address], Transaction{
		Type:      "deposit",
		Amount:    amount,
		Timestamp: time.Now(),
	})

	return nil
}

// Withdraw handles token withdrawals
func (tt *TokenTracker) Withdraw(ctx context.Context, address string, amount uint64) error {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	currentBalance := tt.balances[address]
	if currentBalance < amount {
		return fmt.Errorf("insufficient balance")
	}

	tt.balances[address] -= amount
	tt.history[address] = append(tt.history[address], Transaction{
		Type:      "withdraw",
		Amount:    amount,
		Timestamp: time.Now(),
	})

	return nil
}

// GetBalance returns the current balance for an address
func (tt *TokenTracker) GetBalance(ctx context.Context, address string) (uint64, error) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()

	return tt.balances[address], nil
}

// GetTransactionHistory returns the transaction history for an address
func (tt *TokenTracker) GetTransactionHistory(ctx context.Context, address string) ([]Transaction, error) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()

	return tt.history[address], nil
}
