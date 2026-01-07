# REST API Documentation: JWT-Based Join Flow

This document details the API endpoints for generating join tokens and joining workspaces or boards using JWT tokens. This flow replaces the previous passcode-based join mechanism.

---

## 1. Authentication
All endpoints require a valid JWT Access Token in the request header:
`Authorization: Bearer <access_token>`

---

## 2. Workspaces

### Generate Join Token
Generates a signed JWT token that others can use to join the workspace.
> [!NOTE]
> Only the **creator** of the workspace can generate a join token.

- **URL:** `/api/v1/workspaces/:id/join-token`
- **Method:** `GET`
- **URL Params:** `id=[uint]` (Workspace ID)
- **Response Example:**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

### Join Workspace
Joins a workspace using a valid join token.

- **URL:** `/api/v1/workspaces/join`
- **Method:** `POST`
- **Header:** `Content-Type: application/json`
- **Request Body:**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
- **Response Success (200 OK):**
    ```json
    {
      "message": "Joined workspace successfully"
    }
    ```
- **Possible Errors:**
    - `401 Unauthorized`: Invalid or missing access token.
    - `400 Bad Request`: Missing token in body.
    - `500 Internal Server Error`: Invalid/Expired token, workspace not found, or user already a member.

---

## 3. Boards

### Generate Join Token
Generates a signed JWT token that others can use to join the board.
> [!NOTE]
> Only the **creator** of the board can generate a join token.

- **URL:** `/api/v1/boards/:id/join-token`
- **Method:** `GET`
- **URL Params:** `id=[uint]` (Board ID)
- **Response Example:**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

### Join Board
Joins a board using a valid join token.

- **URL:** `/api/v1/boards/join`
- **Method:** `POST`
- **Header:** `Content-Type: application/json`
- **Request Body:**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
- **Response Success (200 OK):**
    ```json
    {
      "message": "Joined board successfully"
    }
    ```
- **Possible Errors:**
    - `401 Unauthorized`: Invalid or missing access token.
    - `400 Bad Request`: Missing token in body.
    - `500 Internal Server Error`: Invalid/Expired token, board not found, or user already a member.

---

## Summary of URL Patterns
When integrating into the frontend, you might want to create shareable links. For example:
`https://app.putratek.my.id/join?token=eyJhbGciOiJ...`

The frontend should then extract the `token` from the URL and call the `POST /join` endpoint of the respective entity.
