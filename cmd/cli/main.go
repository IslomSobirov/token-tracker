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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

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

	t, err := tracker.NewTokenTracker(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize token tracker: %v", err)
	}
	defer t.Close()

	if err := runCLI(ctx, t); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}

func runCLI(ctx context.Context, t *tracker.TokenTracker) error {
	for {
		fmt.Println("\nToken Tracker")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Deposit")
		fmt.Println("3. Withdraw")
		fmt.Println("4. History")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("> ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var address string
			fmt.Print("Address: ")
			fmt.Scan(&address)
			balance, err := t.GetBalance(ctx, address)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Balance: %d\n", balance)

		case 2:
			var address string
			var amount uint64
			fmt.Print("Address: ")
			fmt.Scan(&address)
			fmt.Print("Amount: ")
			fmt.Scan(&amount)
			if err := t.Deposit(ctx, address, amount); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Deposit successful")

		case 3:
			var address string
			var amount uint64
			fmt.Print("Address: ")
			fmt.Scan(&address)
			fmt.Print("Amount: ")
			fmt.Scan(&amount)
			if err := t.Withdraw(ctx, address, amount); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Withdrawal successful")

		case 4:
			var address string
			fmt.Print("Address: ")
			fmt.Scan(&address)
			history, err := t.GetTransactionHistory(ctx, address)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			for _, tx := range history {
				fmt.Printf("%s: %d (%s)\n",
					tx.Type, tx.Amount, tx.Timestamp.Format("2006-01-02 15:04:05"))
			}

		case 5:
			return nil

		default:
			fmt.Println("Invalid option")
		}
	}
}
