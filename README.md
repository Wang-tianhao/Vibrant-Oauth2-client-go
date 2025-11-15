# Vibrant OAuth2 Client for Go

A simple, thread-safe Go library for authenticating with Vibrant OAuth2 using the client credentials grant type. This library handles token caching and automatic refresh, so you can focus on building your application.

## Features

- Simple one-function API: just call `GetToken()`
- Automatic token caching and refresh
- Thread-safe for concurrent usage
- Configurable via environment variables
- No external dependencies beyond the Go standard library

## Installation

```bash
go get github.com/Wang-tianhao/Vibrant-Oauth2-client-go
```

## Configuration

Set the following environment variables:

```bash
export VIBRANT_CLIENT_ID="your-client-id"
export VIBRANT_CLIENT_SECRET="your-client-secret"
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    vibrant "github.com/Wang-tianhao/Vibrant-Oauth2-client-go"
)

func main() {
    // Create a new client (reads from environment variables)
    client, err := vibrant.NewClient()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Get a token - it will be cached and auto-refreshed
    token, err := client.GetToken()
    if err != nil {
        log.Fatalf("Failed to get token: %v", err)
    }

    fmt.Printf("Access Token: %s\n", token)

    // Subsequent calls will use the cached token if still valid
    token2, err := client.GetToken()
    if err != nil {
        log.Fatalf("Failed to get token: %v", err)
    }

    fmt.Printf("Cached Token: %s\n", token2)
}
```

## How It Works

1. When you call `GetToken()`, the library first checks if there's a cached token
2. If the cached token exists and hasn't expired (with a 60-second buffer), it returns the cached token
3. If there's no cached token or it has expired, the library automatically fetches a new token from the Vibrant OAuth endpoint
4. The new token is cached with its expiration time
5. All operations are thread-safe, so you can safely call `GetToken()` from multiple goroutines

## API Reference

### `NewClient() (*Client, error)`

Creates a new Vibrant OAuth2 client. Reads credentials from environment variables:
- `VIBRANT_CLIENT_ID`
- `VIBRANT_CLIENT_SECRET`

Returns an error if environment variables are not set.

### `GetToken() (string, error)`

Returns a valid access token. This is the main function you'll use. It automatically handles caching and refresh, so you don't need to worry about token expiration.

### `ClearCache()`

Clears the cached token, forcing a new token fetch on the next `GetToken()` call. This is useful for testing or if you need to force a refresh.

## Error Handling

The library returns errors in the following cases:
- Missing environment variables
- Network errors when communicating with the OAuth endpoint
- Invalid responses from the OAuth endpoint
- HTTP errors (non-200 status codes)

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
