package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	sf "github.com/artyconst/fb-storefront-go"
	config "github.com/artyconst/fb-storefront-go/pkg/config"
	customerSDK "github.com/artyconst/fb-storefront-go/pkg/resources/customer"

	env "github.com/artyconst/fb-storefront-go/examples/env"
)

func main() {
	fmt.Println("Fleetbase Storefront - Customer Authentication Test")
	fmt.Println("====================================================")

	envPath := getEnvFilePath()
	if err := env.LoadFromFile(envPath); err != nil {
		log.Printf("Warning: Could not load .env file: %v. Continuing without environment variables.", err)
	}

	apiKey := os.Getenv("STOREFRONT_KEY")
	if apiKey == "" {
		log.Fatal("STOREFRONT_KEY environment variable is required. Please set it in your .env file.")
	}

	client, err := sf.NewStorefront(apiKey,
		config.WithAPIHost(os.Getenv("FLEETBASE_HOST")),
		config.WithLogLevel(config.LevelDebug),
	)
	if err != nil {
		log.Fatal(err)
	}

	customerService := customerSDK.NewCustomerService(client)

	testEmail := os.Getenv("CUSTOMER_EMAIL")
	if testEmail == "" {
		log.Fatal("CUSTOMER_EMAIL environment variable is required. Please set it in your .env file.")
	}

	testPassword := os.Getenv("CUSTOMER_PASSWORD")
	if testPassword == "" {
		log.Fatal("CUSTOMER_PASSWORD environment variable is required. Please set it in your .env file.")
	}

	customerName := os.Getenv("CUSTOMER_NAME")
	if customerName == "" {
		customerName = "Test Customer"
	}

	fmt.Printf("\nUsing email: %s\n", testEmail)
	fmt.Println()

	err = runAuthenticationFlow(customerService, testEmail, testPassword, customerName)
	if err != nil {
		log.Fatalf("Test failed: %v", err)
	}

	fmt.Println("\n✓ Customer authentication test completed successfully!")
}

func runAuthenticationFlow(customerService *customerSDK.CustomerService, email, password, name string) error {
	fmt.Println("Step 1: Attempting customer login...")
	loginResp, loginErr := customerService.Login(context.Background(), customerSDK.LoginRequest{
		Identity: email,
		Password: password,
	})

	if loginErr != nil {
		apiErr := &sf.APIError{}
		if errors.As(loginErr, &apiErr) {
			fmt.Printf("   Login error (code: %s): %s\n", apiErr.Code, apiErr.Message)
		} else {
			fmt.Printf("   Login failed: %v\n", loginErr)
		}

		if strings.Contains(loginErr.Error(), "401") ||
			strings.Contains(loginErr.Error(), "404") ||
			strings.Contains(apiErr.Code, "not_found") ||
			strings.Contains(apiErr.Code, "unauthorized") {
			fmt.Println("\nStep 2: Account not found or unauthorized. Creating new customer...")

			createReq := customerSDK.CustomerCreateRequest{
				Name:     stringPtr(name),
				Identity: email,
				Email:    &email,
				Type:     stringPtr("customer"),
			}

			fmt.Println("   Creating customer with email:", email)
			newCustomer, createErr := customerService.Create(context.Background(), createReq)
			if createErr != nil {
				return fmt.Errorf("failed to create customer: %w", createErr)
			}

			fmt.Printf("   ✓ Customer created successfully!\n")
			fmt.Printf("     ID: %s\n", newCustomer.ID)
			fmt.Printf("     Name: %s\n", derefString(newCustomer.Name))
			fmt.Printf("     Email: %s\n", derefString(newCustomer.Email))

			fmt.Println("\nStep 3: Attempting login with newly created account...")
			time.Sleep(1 * time.Second)

			loginResp2, err := customerService.Login(context.Background(), customerSDK.LoginRequest{
				Identity: email,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("login failed after account creation: %w", err)
			}

			loginResp = loginResp2
		} else {
			return fmt.Errorf("unexpected login error: %w", loginErr)
		}
	}

	fmt.Println("\nStep 4: Login successful!")
	if loginResp.Customer == nil {
		fmt.Printf("   Account successfully authenticated!\n")
	} else {
		fmt.Printf("   Customer ID: %s\n", loginResp.Customer.ID)
		fmt.Printf("   Name: %s\n", derefString(loginResp.Customer.Name))
		fmt.Printf("   Email: %s\n", derefString(loginResp.Customer.Email))
	}
	fmt.Printf("   Token (first 20 chars): %s...\n", loginResp.Token[:min(len(loginResp.Token), 20)])
	if loginResp.ExpiresAt != "" {
		fmt.Printf("   Expires at: %s\n", loginResp.ExpiresAt)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper functions for pointer manipulation
func stringPtr(s string) *string {
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getEnvFilePath returns the appropriate path to .env file based on current working directory.
// This allows examples to work when run from any location (project root, subdirectory, etc.)
func getEnvFilePath() string {
	dir := "examples/.env"
	if wd, err := os.Getwd(); err == nil && strings.Contains(wd, "storefront-go") {
		dir = ".env"
	} else if wd != "" {
		dir = filepath.Join("examples", ".env")
	}

	return dir
}
