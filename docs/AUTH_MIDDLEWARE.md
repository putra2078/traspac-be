# Authentication Middleware Usage Guide

This document explains how the authentication middleware, sliding window, and inactivity timeout work in this application.

## Overview

The authentication system uses a hybrid approach combining JWT (stateless) and Redis (stateful) to provide a "Sliding Window" session. This ensures that users remain logged in as long as they are active, but are automatically logged out after a period of inactivity.

## Components

### 1. Access Token (JWT)
- **Hard Expiration**: 24 Hours (`token_ttl_minute`).
- **Internal Payload**: Contains `user_id` and `email`.
- **Stateless Validation**: The middleware first validates the JWT's signature and its own internal expiration.

### 2. Redis Session
- **Inactivity Timeout**: 15 Minutes (`expires_in_minute`).
- **Storage**: When a user logs in, the `access_token` is stored in Redis as a key.
- **Validation**: The middleware checks if the token exists and is active in Redis.

### 3. Browser Cookie
- **Expiration**: Matches the Inactivity Timeout (15 Minutes).
- **Security**: Set with `HttpOnly` and `Path=/`.

## How the Sliding Window Works

The "Sliding Window" logic is implemented in `internal/middleware/auth.go`:

1.  **Request Arrival**: The middleware extracts the `access_token` from the cookie or header.
2.  **JWT Validation**: Checks if the token is a valid JWT and hasn't exceeded its 24-hour hard limit.
3.  **Redis Check**: Checks if the token is still in Redis (meaning the 15-minute inactivity timer hasn't expired).
4.  **Session Extension**: If valid, the middleware:
    -   Performs `EXPIRE` in Redis to reset the 15-minute timer.
    -   Sends a `Set-Cookie` header to the browser to reset the cookie's 15-minute timer.
5.  **Success**: The request proceeds to the handler.

## Configuration

Settings are found in `config/config.yaml`:

```yaml
jwt:
  secret: "your-secret"
  expires_in_minute: 15    # Inactivity Timeout (Sliding Window)
  token_ttl_minute: 1440   # Hard JWT Expiration (24 Hours)
  refresh_expires_in_days: 7
```

## Middleware Usecase

To protect a route group, use the `AuthMiddleware`:

```go
protected := api.Group("/protected")
protected.Use(middleware.AuthMiddleware(cfg))
{
    protected.GET("/profile", userHandler.GetProfile)
}
```

Every successful request to these routes will automatically extend the user's session.
