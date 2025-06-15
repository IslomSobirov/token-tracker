package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"token-tracker/internal/config"
	"token-tracker/internal/tracker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Parse command line flags
	rpcEndpoint := flag.String("rpc", "https://api.devnet.solana.com", "Solana RPC endpoint")
	wsEndpoint := flag.String("ws", "wss://api.devnet.solana.com", "Solana WebSocket endpoint")
	programID := flag.String("program", "", "Program ID for the token tracker")
	flag.Parse()

	if *programID == "" {
		log.Fatal("Program ID is required")
	}

	cfg := &config.Config{
		RPCEndpoint: *rpcEndpoint,
		WSEndpoint:  *wsEndpoint,
		ProgramID:   *programID,
	}

	// Initialize token tracker
	t, err := tracker.NewTokenTracker(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize token tracker: %v", err)
	}
	defer t.Close()

	// Run the CLI
	if err := runCLI(ctx, t); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}

func runCLI(ctx context.Context, t *tracker.TokenTracker) error {
	for {
		fmt.Println("\nToken Tracker CLI")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Deposit Tokens")
		fmt.Println("3. Withdraw Tokens")
		fmt.Println("4. View Transaction History")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("Select an option: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var address string
			fmt.Print("Enter wallet address: ")
			fmt.Scan(&address)
			balance, err := t.GetBalance(ctx, address)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Balance: %d tokens\n", balance)

		case 2:
			var address string
			var amount uint64
			fmt.Print("Enter wallet address: ")
			fmt.Scan(&address)
			fmt.Print("Enter amount to deposit: ")
			fmt.Scan(&amount)
			if err := t.Deposit(ctx, address, amount); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Deposit successful")

		case 3:
			var address string
			var amount uint64
			fmt.Print("Enter wallet address: ")
			fmt.Scan(&address)
			fmt.Print("Enter amount to withdraw: ")
			fmt.Scan(&amount)
			if err := t.Withdraw(ctx, address, amount); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Withdrawal successful")

		case 4:
			var address string
			fmt.Print("Enter wallet address: ")
			fmt.Scan(&address)
			history, err := t.GetTransactionHistory(ctx, address)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			for _, tx := range history {
				fmt.Printf("Type: %s, Amount: %d, Timestamp: %s\n",
					tx.Type, tx.Amount, tx.Timestamp.Format("2006-01-02 15:04:05"))
			}

		case 5:
			return nil

		default:
			fmt.Println("Invalid option")
		}
	}
}
