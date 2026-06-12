package main

import (
	"context"
	"crypto/rand"
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
	locationSDK "github.com/artyconst/fb-storefront-go/pkg/resources/location"
	productSDK "github.com/artyconst/fb-storefront-go/pkg/resources/product"
	storeSDK "github.com/artyconst/fb-storefront-go/pkg/resources/store"

	env "github.com/artyconst/fb-storefront-go/examples/env"
)

func main() {
	fmt.Println("Fleetbase Storefront Go SDK - End-to-End Shopping Demo")
	fmt.Println("=======================================================")

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
		config.WithUserAgent("MyApp/1.0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	storeService := storeSDK.NewStoreService(client)
	locService := locationSDK.NewService(client)
	cartService := cartSDK.NewCartService(client)
	categoryService := categorySDK.NewCategoryService(client)
	productService := productSDK.NewProductService(client)
	checkoutService := checkoutSDK.NewCheckoutService(client)

	ctx := context.Background()

	// === Phase 1: Discover Store & Products ===
	fmt.Println("\n=== Phase 1: Discover Store & Products ===")

	about, err := storeService.About(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Store: %s (ID: %s)\n", about.Name, about.ID)

	categories, err := categoryService.List(ctx, categorySDK.WithLimit(10))
	if err != nil {
		fmt.Printf("   Note: List categories failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d categories:\n", len(categories))
		for _, cat := range categories[:min(len(categories), 3)] {
			fmt.Printf("     - %s (ID: %s)\n", cat.Name, cat.ID)
		}
		if len(categories) > 3 {
			fmt.Printf("     ... and %d more\n", len(categories)-3)
		}
	}

	products, err := productService.List(ctx, productSDK.WithLimit(10))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to list products: %w", err))
	}
	fmt.Printf("   Found %d products:\n", len(products))
	for i, prod := range products[:min(len(products), 5)] {
		variantsInfo := ""
		if len(prod.Variants) > 0 {
			variantNames := make([]string, 0, len(prod.Variants))
			for _, v := range prod.Variants {
				opts := make([]string, 0, len(v.Options))
				for _, o := range v.Options {
					opts = append(opts, fmt.Sprintf("%s=%s", o.Name, o.Value))
				}
				variantNames = append(variantNames, fmt.Sprintf("[%s: %v]", v.Name, strings.Join(opts, ", ")))
			}
			variantsInfo = " | Variants: " + strings.Join(variantNames, "; ")
		}
		fmt.Printf("     %d. %s (ID: %s) - $%d%s\n", i+1, prod.Name, prod.ID, prod.Price, variantsInfo)
	}

	if len(products) == 0 {
		log.Fatal("No products found in the store. Cannot continue with demo.")
	}

	// === Phase 2: Create Cart & Add Items ===
	fmt.Println("\n=== Phase 2: Create Cart & Add Items ===")

	cartID := generateUUID()
	fmt.Printf("   Generated cart ID: %s\n", cartID)

	// Create the cart implicitly via Get (auto-creates if it doesn't exist)
	cart, err := cartService.Get(ctx, cartID)
	if err != nil {
		fmt.Printf("   Note: Get cart failed: %v\n", err)
	} else {
		fmt.Printf("   Cart created - Status: %s, Items: %d\n", cart.Status, len(cart.Items))
	}

	// --- Add item 1: Simple add without variants ---
	if len(products) > 0 {
		simpleProduct := products[0]
		fmt.Printf("\n   Adding simple product (no variants): %s x1\n", simpleProduct.Name)
		cart, err = cartService.AddItem(ctx, cartID, simpleProduct.ID, 1)
		if err != nil {
			fmt.Printf("   Note: AddItem failed: %v\n", err)
		} else {
			fmt.Printf("   Cart now has %d item(s), Total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)
		}
	}

	// --- Add item 2: Product with variants (if available) ---
	variantProduct := findProductWithVariants(products)
	if variantProduct != nil {
		fmt.Printf("\n   Adding product with variants: %s x1\n", variantProduct.Name)

		// Build variant options from the first variant's options
		opts := make([]cartSDK.AddItemOption, 0)
		for _, opt := range variantProduct.Variants[0].Options {
			fmt.Printf("     Variant option: %s = %s\n", opt.Name, opt.Value)
			opts = append(opts, cartSDK.WithVariant(opt.Name, opt.Value))
		}

		cart, err = cartService.AddItem(ctx, cartID, variantProduct.ID, 1, opts...)
		if err != nil {
			fmt.Printf("   Note: AddItem with variants failed: %v\n", err)
		} else {
			fmt.Printf("   Cart now has %d item(s), Total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)
		}
	}

	// --- Add item 3: With scheduled time ---
	if len(products) > 1 {
		scheduledProduct := products[1]
		scheduledTime := "2024-12-25T10:00:00Z"
		fmt.Printf("\n   Adding product with scheduled delivery: %s x2 (scheduled for %s)\n", scheduledProduct.Name, scheduledTime)

		cart, err = cartService.AddItem(ctx, cartID, scheduledProduct.ID, 2,
			cartSDK.WithScheduledAt(scheduledTime),
		)
		if err != nil {
			fmt.Printf("   Note: AddItem with schedule failed: %v\n", err)
		} else {
			fmt.Printf("   Cart now has %d item(s), Total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)
		}
	}

	// --- Add item 4: With store location (if locations available) ---
	if len(products) > 2 {
		locations, err := locService.List(ctx, about.ID)
		if err != nil || len(locations) == 0 {
			fmt.Printf("\n   Note: No store locations available to demonstrate WithStoreLocation\n")
		} else {
			locationProduct := products[2]
			storeLocID := locations[0].ID
			fmt.Printf("\n   Adding product with store location: %s x1 (location: %s)\n", locationProduct.Name, locations[0].Name)

			cart, err = cartService.AddItem(ctx, cartID, locationProduct.ID, 1,
				cartSDK.WithStoreLocation(storeLocID),
			)
			if err != nil {
				fmt.Printf("   Note: AddItem with store location failed: %v\n", err)
			} else {
				fmt.Printf("   Cart now has %d item(s), Total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)
			}
		}
	}

	// === Phase 3: Update Cart Item ===
	fmt.Println("\n=== Phase 3: Update Cart Item ===")

	cart, err = cartService.Get(ctx, cartID)
	if err != nil {
		fmt.Printf("   Note: Get cart failed: %v\n", err)
	} else if len(cart.Items) > 0 {
		firstItem := cart.Items[0]
		newQuantity := firstItem.Quantity + 1
		fmt.Printf("   Current items in cart:\n")
		for _, item := range cart.Items {
			fmt.Printf("     - %s (ID: %s) x%d @ $%d = $%d\n", item.Name, item.ID, item.Quantity, item.Price, item.Total)
		}

		fmt.Printf("\n   Updating first item '%s' quantity from %d to %d\n", firstItem.Name, firstItem.Quantity, newQuantity)

		cart, err = cartService.UpdateItem(ctx, cartID, firstItem.ID, newQuantity)
		if err != nil {
			fmt.Printf("   Note: UpdateItem failed: %v\n", err)
		} else {
			fmt.Printf("   Cart updated - Total items: %d, New total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)
			for _, item := range cart.Items {
				fmt.Printf("     - %s x%d @ $%d = $%d\n", item.Name, item.Quantity, item.Price, item.Total)
			}
		}
	}

	// === Phase 4: Checkout Flow ===
	fmt.Println("\n=== Phase 4: Checkout Flow ===")

	cart, err = cartService.Get(ctx, cartID)
	if err != nil {
		fmt.Printf("   Note: Get cart failed before checkout: %v\n", err)
	} else if len(cart.Items) == 0 {
		fmt.Println("   Cart is empty, skipping checkout.")
	} else {
		fmt.Printf("   Proceeding to checkout with %d item(s), Total: $%d (%s)\n", len(cart.Items), cart.TotalAmount, cart.Currency)

		checkoutReq := checkoutSDK.CreateCheckoutRequest{
			CustomerEmail: "demo@example.com",
			ShippingAddress: &checkoutSDK.Address{
				FirstName:    "Demo",
				LastName:     "User",
				AddressLine1: "123 Main St",
				City:         "Manila",
				PostalCode:   "1000",
				Country:      "PH",
			},
		}

		checkout, err := checkoutService.Create(ctx, cartID, checkoutReq)
		if err != nil {
			fmt.Printf("   Note: Create checkout failed: %v\n", err)
		} else {
			fmt.Printf("   Checkout created - ID: %s, Status: %s, Amount: $%d (%s)\n",
				checkout.ID, checkout.Status, checkout.Amount, checkout.Currency)

			if checkout.ShippingAddress != nil {
				fmt.Printf("     Shipping to: %s %s, %s\n",
					checkout.ShippingAddress.FirstName,
					checkout.ShippingAddress.LastName,
					checkout.ShippingAddress.City)
			}

			// Capture the checkout as an order with notes
			notes := "Demo order from Go SDK"
			captured, err := checkoutService.CaptureCheckout(ctx, "", checkoutSDK.WithNotes(notes))
			if err != nil {
				fmt.Printf("   Note: Capture checkout failed: %v\n", err)
			} else {
				fmt.Printf("   Checkout captured - ID: %s, Status: %s\n", captured.ID, captured.Status)
			}
		}
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("This demo demonstrated the full shopping flow:")
	fmt.Println("  1. Browsing store info, categories, and products (with variant details)")
	fmt.Println("  2. Creating a cart and adding items with various options:")
	fmt.Println("     - Simple add item (no variants)")
	fmt.Println("     - Add item with product variants")
	fmt.Println("     - Add item with scheduled delivery time")
	fmt.Println("     - Add item with store location")
	fmt.Println("  3. Updating a cart item quantity")
	fmt.Println("  4. Creating and capturing a checkout session")
}

// findProductWithVariants returns the first product that has variants defined, or nil.
func findProductWithVariants(products []*productSDK.Product) *productSDK.Product {
	for _, p := range products {
		if len(p.Variants) > 0 {
			return p
		}
	}
	return nil
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

// generateUUID returns a random UUID v4 string using crypto/rand.
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}


