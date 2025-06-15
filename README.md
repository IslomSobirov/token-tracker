# Token Tracker

A Go-based token tracking application that interacts with the Solana blockchain. This application provides a CLI interface for tracking token balances and managing deposits/withdrawals.

## Features

- Token balance tracking
- Deposit and withdrawal functionality
- Transaction history tracking
- Real-time balance updates
- Graceful shutdown handling
- Thread-safe operations

## Project Structure

```
.
├── cmd/
│   └── cli/           # CLI application
├── internal/
│   ├── config/        # Configuration management
│   └── tracker/       # Core token tracking functionality
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21 or later
- Solana CLI tools (for deployment)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd token-tracker
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o token-tracker ./cmd/cli
```

## Usage

Run the CLI application:
```bash
./token-tracker -program <your-program-id>
```

Available flags:
- `-rpc`: Solana RPC endpoint (default: https://api.devnet.solana.com)
- `-ws`: Solana WebSocket endpoint (default: wss://api.devnet.solana.com)
- `-program`: Program ID for the token tracker (required)

### CLI Commands

The application provides an interactive CLI with the following options:

1. Check Balance
   - Enter a wallet address to view its current token balance

2. Deposit Tokens
   - Enter a wallet address and amount to deposit tokens

3. Withdraw Tokens
   - Enter a wallet address and amount to withdraw tokens

4. View Transaction History
   - Enter a wallet address to view its transaction history

5. Exit
   - Gracefully shut down the application

## Deployment to Testnet

### 1. Install Solana CLI Tools

On macOS:
```bash
sh -c "$(curl -sSfL https://release.solana.com/v1.17.9/install)"
```

On Linux:
```bash
sh -c "$(curl -sSfL https://release.solana.com/v1.17.9/install)"
```

On Windows:
Download the installer from https://docs.solana.com/cli/install-solana-cli-tools

### 2. Configure for Testnet

```bash
solana config set --url https://api.testnet.solana.com
```

### 3. Create a New Wallet

```bash
solana-keygen new
```

### 4. Get Your Wallet Address

```bash
solana address
```

### 5. Get Test SOL

Request test SOL from the faucet:
```bash
solana airdrop 2
```

### 6. Build and Deploy

1. Build the application:
```bash
go build -o token-tracker ./cmd/cli
```

2. Deploy to testnet:
```bash
solana program deploy token-tracker
```

3. Note the program ID from the deployment output and use it to run the application:
```bash
./token-tracker -program <deployed-program-id>
```

### 7. Verify Deployment

Check your program's status:
```bash
solana program show <deployed-program-id>
```

## Development

### Adding New Features

1. Core functionality should be added to the `internal/tracker` package
2. Configuration changes should be made in the `internal/config` package
3. CLI interface modifications should be made in `cmd/cli`

### Testing

Run the tests:
```bash
go test ./...
```

## Deployment

The application can be deployed to different Solana networks by specifying the appropriate RPC endpoints:

- Devnet: https://api.devnet.solana.com
- Testnet: https://api.testnet.solana.com
- Mainnet: https://api.mainnet-beta.solana.com

## Security Considerations

- Always use secure RPC endpoints
- Implement proper transaction signing
- Add rate limiting for production use
- Implement proper error handling
- Add logging and monitoring
- Keep your private keys secure
- Use environment variables for sensitive data
- Implement proper access controls

## Troubleshooting

### Common Issues

1. Insufficient SOL for deployment
   - Solution: Request more test SOL using `solana airdrop 2`

2. Program deployment fails
   - Check your SOL balance
   - Verify network connectivity
   - Ensure you're using the correct network (testnet)

3. Transaction failures
   - Verify account balances
   - Check transaction signatures
   - Ensure proper program permissions

## License

MIT 