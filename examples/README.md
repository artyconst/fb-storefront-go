# Fleetbase Storefront Go SDK Examples

This directory contains example applications demonstrating how to use the Fleetbase Storefront Go SDK.

## Setup

### 1. Copy Environment Variables

Copy the `.env.example` file to `.env` and fill in your actual credentials:

```bash
cp .env.example .env
```

Edit `.env` with your values:

```bash
STOREFRONT_KEY=YOUR_API_KEY
CUSTOMER_EMAIL=<your-email>
CUSTOMER_PASSWORD=<your-password>
CUSTOMER_NAME=<your-customer-name>
FLEETBASE_HOST=https://api.storefront.fleetbase.io/v1
```

### 2. Run Examples

#### Demo Example

Demonstrates all major SDK features including products, carts, orders, and checkout:

```bash
cd demo
go run main.go
```

#### Customer Authentication Example

Demonstrates customer authentication flows (login, registration, token handling):

```bash
cd customer_auth
go run main.go
```

## How It Works

### Environment Configuration Without Dependencies

The examples use a custom `.env` loader (`examples/env/env.go`) that parses environment files using only Go's standard library. This approach:

- **No third-party dependencies**: The `godotenv` package is not required, keeping the SDK lightweight
- **Shared code**: Both examples import from the shared `examples/env` package to avoid duplication
- **Simple structure**: All examples are part of the main SDK module and can import each other directly

## Requirements

- Go 1.25.5 or later
- A valid Fleetbase Storefront API key (can be obtained from your Fleetbase dashboard)

## Notes

- The examples demonstrate common usage patterns and can be run directly with `go run`
- All sensitive credentials should be stored in `.env` files, which are ignored by version control (see .gitignore)
- For production use, always set environment variables directly instead of using `.env` files
