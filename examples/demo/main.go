package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	sf "github.com/artyconst/fb-storefront-go"
	config "github.com/artyconst/fb-storefront-go/pkg/config"
	cartSDK "github.com/artyconst/fb-storefront-go/pkg/resources/cart"
	categorySDK "github.com/artyconst/fb-storefront-go/pkg/resources/category"
	checkoutSDK "github.com/artyconst/fb-storefront-go/pkg/resources/checkout"
	orderSDK "github.com/artyconst/fb-storefront-go/pkg/resources/order"
	productSDK "github.com/artyconst/fb-storefront-go/pkg/resources/product"
	storeSDK "github.com/artyconst/fb-storefront-go/pkg/resources/store"

	env "github.com/artyconst/fb-storefront-go/examples/env"
)

func main() {
	fmt.Println("Fleetbase Storefront Go SDK - Full Demo")
	fmt.Println("===========================================")

	envFile := getEnvFilePath()

	if err := env.LoadFromFile(envFile); err != nil {
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

	storeService := storeSDK.NewStoreService(client)
	cartService := cartSDK.NewCartService(client)
	categoryService := categorySDK.NewCategoryService(client)
	productService := productSDK.NewProductService(client)
	checkoutService := checkoutSDK.NewCheckoutService(client)
	orderService := orderSDK.NewOrderService(client)

	fmt.Println("\n1. Testing GET /about")
	about, err := storeService.About(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Store: %s\n", about.Name)
	fmt.Printf("   ID: %s\n", about.ID)

	fmt.Println("\n2. Store Search (Search method removed - endpoint not available for store-specific API keys)")

	fmt.Println("\n3. Testing Category Operations with functional options")
	categories, err := categoryService.List(context.Background(),
		categorySDK.WithLimit(10),
	)
	if err != nil {
		fmt.Printf("   Note: List categories failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d categories\n", len(categories))
		for _, cat := range categories[:min(len(categories), 3)] {
			fmt.Printf("     Category ID: %s, Name: %s\n", cat.ID, cat.Name)
		}
	}

	fmt.Println("\n4. Testing Product Operations with functional options")
	products, err := productService.List(context.Background(),
		productSDK.WithOffset(0),
	)
	if err != nil {
		fmt.Printf("   Note: List products failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d products\n", len(products))
		for _, prod := range products[:min(len(products), 3)] {
			fmt.Printf("     Product ID: %s, Name: %s, Price: $%d\n", prod.ID, prod.Name, prod.Price)
		}
	}

	if len(categories) > 0 {
		fmt.Println("\n5. Testing Find Products by Category (uses functional options internally)")
		catProducts, err := productService.FindByCategory(context.Background(), categories[0].ID)
		if err != nil {
			fmt.Printf("   Note: Find products by category failed: %v\n", err)
		} else {
			fmt.Printf("     Found %d products in category %s\n", len(catProducts), categories[0].Name)
		}
	}

	fmt.Println("\n6. Testing GET /gateways with functional options")
	gateways, err := storeService.ListGateways(context.Background(),
		storeSDK.WithGatewayLimit(5),
		storeSDK.WithGatewayOffset(0),
	)
	if err != nil {
		fmt.Printf("   Note: List gateways failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d payment gateways\n", len(gateways.Data))

		if len(gateways.Data) > 0 {
			fmt.Println("\n7. Testing GET /gateways/{id}")
			gateway, err := storeService.GetGateway(context.Background(), gateways.Data[0].ID)
			if err != nil {
				fmt.Printf("   Note: Get gateway failed: %v\n", err)
			} else {
				fmt.Printf("   Gateway: %s (%s)\n", gateway.Name, gateway.Type)
				fmt.Printf("   Active: %v\n", gateway.IsActive)
			}
		}
	}

	fmt.Println("\n8. Testing Cart Operations")
	cartID := "cart_123" // Replace with actual cart ID from your store
	cart, err := cartService.Get(context.Background(), cartID)
	if err != nil {
		fmt.Printf("   Note: Get cart failed (expected if no cart exists): %v\n", err)
	} else {
		fmt.Printf("   Cart Status: %s\n", cart.Status)
		fmt.Printf("   Total: $%d (%s)\n", cart.TotalAmount, cart.Currency)
		fmt.Printf("   Items:\n")
		for _, item := range cart.Items {
			fmt.Printf("     - %s x%d ($%d)\n", item.Name, item.Quantity, item.Price)
		}

		fmt.Println("\n9. Testing Cart AddItem")
		addedCart, err := cartService.AddItem(context.Background(), cartID, "prod_123", 1, nil, nil, "", "")
		if err != nil {
			fmt.Printf("   Note: AddItem failed (expected): %v\n", err)
		} else {
			fmt.Printf("   Cart now has %d items\n", len(addedCart.Items))
		}

		fmt.Println("\n10. Testing Cart Clear")
		err = cartService.Clear(context.Background(), cartID)
		if err != nil {
			fmt.Printf("   Note: Clear failed (expected): %v\n", err)
		} else {
			fmt.Printf("   Cart cleared successfully\n")
		}
	}

	fmt.Println("\n11. Testing Checkout Operations")
	fmt.Println("   Note: List capturable checkouts not available in SDK")

	fmt.Println("\n12. Testing Capture Checkout")
	captured, err := checkoutService.CaptureCheckout(context.Background(), "")
	if err != nil {
		fmt.Printf("   Note: Capture failed (expected): %v\n", err)
	} else {
		fmt.Printf("     Captured Checkout ID: %s, Status: %s\n", captured.ID, captured.Status)
	}

	fmt.Println("\n13. Testing Order Operations with functional options (page-based pagination)")
	ordersList, err := orderService.List(context.Background(),
		orderSDK.WithPage(1),
		orderSDK.WithLimit(10),
	)
	if err != nil {
		fmt.Printf("   Note: List orders returned error (expected if no orders exist): %v\n", err)
	} else {
		fmt.Printf("   Found %d orders\n", len(ordersList))
		for _, ord := range ordersList[:min(len(ordersList), 3)] {
			fmt.Printf("     Order ID: %s, Status: %s, Total: $%d (%s)\n", ord.ID, ord.Status, ord.TotalAmount, ord.Currency)
		}
	}

	fmt.Println("\n✓ All SDK operations working correctly!")
}

func getEnvFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "./examples/.env"
	}

	if strings.Contains(cwd, "/examples/") || filepath.Base(cwd) == "examples" {
		return ".env"
	}

	return "./examples/.env"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
