package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sf "github.com/artyconst/fb-storefront-go"
	config "github.com/artyconst/fb-storefront-go/pkg/config"
	cartSDK "github.com/artyconst/fb-storefront-go/pkg/resources/cart"
	checkoutSDK "github.com/artyconst/fb-storefront-go/pkg/resources/checkout"
	customerSDK "github.com/artyconst/fb-storefront-go/pkg/resources/customer"
	foodtruckSDK "github.com/artyconst/fb-storefront-go/pkg/resources/foodtruck"
	locationSDK "github.com/artyconst/fb-storefront-go/pkg/resources/location"
	orderSDK "github.com/artyconst/fb-storefront-go/pkg/resources/order"
	productSDK "github.com/artyconst/fb-storefront-go/pkg/resources/product"
	reviewSDK "github.com/artyconst/fb-storefront-go/pkg/resources/review"

	env "github.com/artyconst/fb-storefront-go/examples/env"
)

func main() {
	fmt.Println("Fleetbase Storefront Go SDK - New Features Demo")
	fmt.Println("=================================================")

	// Load environment variables from .env file
	if err := env.LoadFromFile("examples/.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v. Continuing without environment variables.", err)
	}

	apiKey := os.Getenv("STOREFRONT_KEY")
	if apiKey == "" {
		log.Fatal("STOREFRONT_KEY not set in environment or examples/.env")
	}

	host := os.Getenv("FLEETBASE_HOST")
	if host == "" {
		host = "https://api.storefront.fleetbase.io/v1"
	}

	client, err := sf.NewStorefront(apiKey,
		config.WithAPIHost(host),
		config.WithLogLevel(config.LevelDebug),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// --- 1. Food Trucks ---
	demoFoodTrucks(ctx, client)

	// --- 2. Store Locations ---
	demoStoreLocations(ctx, client)

	// --- 3. Cart Mutations (AddProduct, UpdateLineItem, RemoveLineItem, EmptyCart, DeleteCart) ---
	demoCartMutations(ctx, client)

	// --- 4. Checkout New Methods (Initialize, Status, CaptureQPay) ---
	demoCheckoutNewMethods(ctx, client)

	// --- 5. Review CRUD ---
	demoReviewCRUD(ctx, client)

	// --- 6. Order Fulfillment (MarkPickedUp, GenerateReceipt) ---
	demoOrderFulfillment(ctx, client)

	// --- 7. Customer Advanced Features (OAuth, Stripe, Phone Verification, Account Closure) ---
	demoCustomerAdvanced(ctx, client)

	fmt.Println("\nAll new features demo completed!")
}

func demoFoodTrucks(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Food Trucks ===")

	service := foodtruckSDK.NewService(client)

	// List all food trucks
	trucks, err := service.List(ctx)
	if err != nil {
		var apiErr *sf.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API error listing food trucks: %s\n", apiErr.Message)
			return
		}
		log.Fatalf("Failed to list food trucks: %v", err)
	}

	fmt.Printf("Found %d food truck(s)\n", len(trucks))
	for i, truck := range trucks {
		if i >= min(len(trucks), 3) {
			break
		}
		name := "Unknown"
		if truck.Vehicle != nil && truck.Vehicle.Name != "" {
			name = truck.Vehicle.Name
		}
		fmt.Printf("  - %s (ID: %s)\n", name, truck.ID)
	}

	// Get a specific food truck if we have any
	if len(trucks) > 0 {
		truck, err := service.Get(ctx, trucks[0].ID)
		if err != nil {
			fmt.Printf("Failed to get food truck: %v\n", err)
			return
		}
		name := "Unknown"
		if truck.Vehicle != nil && truck.Vehicle.Name != "" {
			name = truck.Vehicle.Name
		}
		fmt.Printf("Food truck details: %s — %d catalogs, status=%s\n", name, len(truck.Catalogs), truck.Status)
	}
}

func demoStoreLocations(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Store Locations ===")

	service := locationSDK.NewService(client)

	locations, err := service.ListLocations(ctx)
	if err != nil {
		var apiErr *sf.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API error listing store locations: %s\n", apiErr.Message)
			return
		}
		log.Fatalf("Failed to list store locations: %v", err)
	}

	fmt.Printf("Found %d store location(s)\n", len(locations))
	for i, loc := range locations {
		if i >= min(len(locations), 3) {
			break
		}
		fmt.Printf("  - %s (ID: %s)\n", loc.Name, loc.ID)
	}
}

func demoCartMutations(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Cart Mutations ===")

	service := cartSDK.NewCartService(client)

	// Create a new cart (implicitly via Get - auto-creates if it doesn't exist)
	cartID := fmt.Sprintf("cart_%d", time.Now().Unix())
	cart, err := service.Get(ctx, cartID)
	if err != nil {
		fmt.Printf("Failed to create cart: %v\n", err)
		return
	}
	fmt.Printf("Created cart: %s (status: %s)\n", cart.ID, cart.Status)

	// Add a product using the new AddProduct method (functional options instead of quantity param)
	productID := "prod_01J..." // Replace with actual product ID
	addedCart, err := service.AddProduct(ctx, cartID, productID,
		cartSDK.WithQuantity(2),
	)
	if err != nil {
		fmt.Printf("AddProduct failed (expected if product doesn't exist): %v\n", err)
	} else {
		fmt.Printf("Added product to cart: %d items\n", len(addedCart.Items))
	}

	// Update a line item using the new UpdateLineItem method (functional options instead of quantity param)
	if len(cart.Items) > 0 {
		lineItemID := cart.Items[0].ID
		_, err = service.UpdateLineItem(ctx, cartID, lineItemID,
			cartSDK.WithQuantityForUpdate(5),
		)
		if err != nil {
			fmt.Printf("UpdateLineItem failed: %v\n", err)
		} else {
			fmt.Println("Updated line item quantity to 5")
		}

		// Remove a line item using the new RemoveLineItem method name
		_, err = service.RemoveLineItem(ctx, cartID, lineItemID)
		if err != nil {
			fmt.Printf("RemoveLineItem failed: %v\n", err)
		} else {
			fmt.Println("Removed line item from cart")
		}
	}

	// Empty the cart using the new EmptyCart method (returns *Cart instead of error)
	emptyCart, err := service.EmptyCart(ctx, cartID)
	if err != nil {
		fmt.Printf("EmptyCart failed: %v\n", err)
	} else {
		fmt.Printf("Emptied cart: %d items remaining\n", len(emptyCart.Items))
	}

	// Delete the cart using the new DeleteCart method
	err = service.DeleteCart(ctx)
	if err != nil {
		fmt.Printf("DeleteCart failed: %v\n", err)
	} else {
		fmt.Println("Deleted cart")
	}
}

func demoCheckoutNewMethods(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Checkout New Methods ===")

	service := checkoutSDK.NewCheckoutService(client)

	// Initialize checkout preview using the new Initialize method (/before endpoint)
	preview, err := service.Initialize(ctx,
		checkoutSDK.WithCart("cart_01J..."), // Replace with actual cart ID
		checkoutSDK.WithCash(),
	)
	if err != nil {
		fmt.Printf("Initialize failed (expected if no valid cart): %v\n", err)
	} else {
		fmt.Printf("Checkout preview: token=%s\n", preview.Token)
	}

	// Get checkout status using the new Status method
	status, err := service.Status(ctx,
		checkoutSDK.WithToken(""), // Replace with actual token
	)
	if err != nil {
		fmt.Printf("Status failed (expected without valid token): %v\n", err)
	} else {
		fmt.Printf("Checkout status: %s\n", status.Status)
	}

	// Capture QPay using the new CaptureQPay method (for QPay callback flow)
	order, err := service.CaptureQPay(ctx, "checkout_01J...", false, nil) // Replace with actual checkout ID
	if err != nil {
		fmt.Printf("CaptureQPay failed (expected if checkout doesn't exist): %v\n", err)
	} else {
		fmt.Printf("Captured QPay checkout: order %s (status: %s)\n", order.ID, order.Status)
	}
}

func demoReviewCRUD(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Review CRUD ===")

	service := reviewSDK.NewReviewService(client)

	// List reviews
	reviews, err := service.List(ctx,
		reviewSDK.WithReviewLimit(5),
	)
	if err != nil {
		var apiErr *sf.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API error listing reviews: %s\n", apiErr.Message)
			return
		}
		log.Fatalf("Failed to list reviews: %v", err)
	}

	fmt.Printf("Found %d review(s)\n", len(reviews))
	for i, r := range reviews {
		if i >= min(len(reviews), 3) {
			break
		}
		fmt.Printf("  - Rating: %d — %s\n", r.Rating, r.Content)
	}

	// Create a review using the new Create method
	newReview, err := service.Create(ctx, "Great product!", 5)
	if err != nil {
		fmt.Printf("Create review failed: %v\n", err)
	} else {
		fmt.Printf("Created review: %s (rating: %d)\n", newReview.ID, newReview.Rating)

		// Get the review by ID using the new Get method
		gotReview, err := service.Get(ctx, newReview.ID)
		if err != nil {
			fmt.Printf("Get review failed: %v\n", err)
		} else {
			fmt.Printf("Retrieved review: rating=%d content=%s\n", gotReview.Rating, gotReview.Content)

			// Delete the review using the new Delete method
			err = service.Delete(ctx, newReview.ID)
			if err != nil {
				fmt.Printf("Delete review failed: %v\n", err)
			} else {
				fmt.Println("Deleted review")
			}
		}
	}
}

func demoOrderFulfillment(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Order Fulfillment ===")

	service := orderSDK.NewOrderService(client)

	orderID := "order_01J..." // Replace with actual order ID

	// Mark order as picked up using the new MarkPickedUp method
	err := service.MarkPickedUp(ctx, orderID)
	if err != nil {
		fmt.Printf("MarkPickedUp failed (expected if order doesn't exist): %v\n", err)
	} else {
		fmt.Println("Order marked as picked up")
	}

	// Generate receipt using the new GenerateReceipt method
	receipt, err := service.GenerateReceipt(ctx, orderID,
		orderSDK.WithEbarimtReceiverType("CITIZEN"),
	)
	if err != nil {
		fmt.Printf("GenerateReceipt failed: %v\n", err)
	} else {
		fmt.Printf("Generated receipt for order %s\n", receipt.OrderID)
	}
}

func demoCustomerAdvanced(ctx context.Context, client *sf.StorefrontClient) {
	fmt.Println("\n=== Customer Advanced Features ===")

	service := customerSDK.NewCustomerService(client)

	// Login with Facebook using the new LoginWithFacebook method (OAuth)
	customer, err := service.LoginWithFacebook(ctx, "facebook_user_123", nil, nil)
	if err != nil {
		fmt.Printf("LoginWithFacebook failed (expected): %v\n", err)
	} else {
		fmt.Printf("Logged in with Facebook: customer %s\n", customer.ID)
	}

	// Get Stripe ephemeral key using the new GetStripeEphemeralKey method
	key, err := service.GetStripeEphemeralKey(ctx)
	if err != nil {
		fmt.Printf("GetStripeEphemeralKey failed (expected without auth): %v\n", err)
	} else {
		fmt.Printf("Got Stripe ephemeral key: %s...\n", key[:min(len(key), 20)])
	}

	// Create Stripe setup intent using the new CreateStripeSetupIntent method
	setupIntent, err := service.CreateStripeSetupIntent(ctx)
	if err != nil {
		fmt.Printf("CreateStripeSetupIntent failed (expected without auth): %v\n", err)
	} else {
		fmt.Printf("Created Stripe setup intent: client_secret=%s...\n", setupIntent.ClientSecret[:min(len(setupIntent.ClientSecret), 20)])
	}

	// Request phone verification using the new RequestPhoneVerification method
	err = service.RequestPhoneVerification(ctx, "+1234567890")
	if err != nil {
		fmt.Printf("RequestPhoneVerification failed: %v\n", err)
	} else {
		fmt.Println("Phone verification requested")
	}

	// Verify phone number using the new VerifyPhoneNumber method
	err = service.VerifyPhoneNumber(ctx, "123456", "+1234567890")
	if err != nil {
		fmt.Printf("VerifyPhoneNumber failed: %v\n", err)
	} else {
		fmt.Println("Phone number verified")
	}

	// Initiate account closure using the new InitiateAccountClosure method
	err = service.InitiateAccountClosure(ctx)
	if err != nil {
		fmt.Printf("InitiateAccountClosure failed (expected without auth): %v\n", err)
	} else {
		fmt.Println("Account closure initiated")
	}

	// Confirm account closure using the new ConfirmAccountClosure method
	err = service.ConfirmAccountClosure(ctx, "closure_code_123")
	if err != nil {
		fmt.Printf("ConfirmAccountClosure failed: %v\n", err)
	} else {
		fmt.Println("Account closure confirmed")
	}

	// Update customer using the new Update method (functional options)
	updatedCustomer, err := service.Update(ctx, "cust_01J...", // Replace with actual customer ID
		customerSDK.WithName("John Doe"),
		customerSDK.WithEmail("john@example.com"),
	)
	if err != nil {
		fmt.Printf("Update customer failed: %v\n", err)
	} else {
		fmt.Printf("Updated customer: %s (ID: %s)\n", updatedCustomer.Name, updatedCustomer.ID)
	}

	// --- Product Write Operations ---
	fmt.Println("\n=== Product Write Operations ===")
	productService := productSDK.NewProductService(client)

	// Create a new product using the new Create method (functional options)
	newProduct, err := productService.Create(ctx,
		productSDK.WithProductName("Test Product"),
		productSDK.WithPrice(1999), // $19.99 in cents
	)
	if err != nil {
		fmt.Printf("Create product failed: %v\n", err)
	} else {
		fmt.Printf("Created product: %s (ID: %s)\n", newProduct.Name, newProduct.ID)

		// Update the product using the new Update method (functional options)
		updatedProduct, err := productService.Update(ctx, newProduct.ID,
			productSDK.WithUpdatedPrice(2499), // $24.99 in cents
		)
		if err != nil {
			fmt.Printf("Update product failed: %v\n", err)
		} else {
			fmt.Printf("Updated product price to: %d\n", updatedProduct.Price)
		}
	}
}
