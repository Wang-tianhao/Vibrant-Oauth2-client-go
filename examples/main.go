package main

import (
	"fmt"
	"log"
	"time"

	vibrant "github.com/Wang-tianhao/Vibrant-Oauth2-client-go"
)

func main() {
	// Create a new Vibrant OAuth2 client
	// It will read VIBRANT_CLIENT_ID and VIBRANT_CLIENT_SECRET from environment
	client, err := vibrant.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Vibrant OAuth2 Client Example")
	fmt.Println("==============================")

	// First call - will fetch a new token
	fmt.Println("\n1. Fetching first token...")
	token1, err := client.GetToken()
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("   Token obtained: %s...\n", token1[:20])

	// Second call - should return cached token
	fmt.Println("\n2. Getting token again (should be cached)...")
	token2, err := client.GetToken()
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("   Token obtained: %s...\n", token2[:20])

	if token1 == token2 {
		fmt.Println("   ✓ Token was returned from cache")
	}

	// Demonstrate concurrent access
	fmt.Println("\n3. Testing concurrent access...")
	for i := 0; i < 5; i++ {
		go func(id int) {
			token, err := client.GetToken()
			if err != nil {
				log.Printf("Goroutine %d failed: %v", id, err)
				return
			}
			fmt.Printf("   Goroutine %d got token: %s...\n", id, token[:20])
		}(i)
	}

	// Wait for goroutines to complete
	time.Sleep(2 * time.Second)

	// Clear cache to force refresh
	fmt.Println("\n4. Clearing cache and fetching new token...")
	client.ClearCache()
	token3, err := client.GetToken()
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("   New token obtained: %s...\n", token3[:20])

	fmt.Println("\n✓ Example completed successfully!")
}
