# Boards API Guide

This document outlines the REST API endpoints for managing Boards.

## Create Board

Creates a new board within a specific workspace.

- **URL**: `/api/v1/boards/`
- **Method**: `POST`
- **Authentication**: Required (`Authorization: Bearer <token>`)

### Request Headers

| Header | Value | Description |
| :--- | :--- | :--- |
| `Content-Type` | `application/json` | |
| `Authorization` | `Bearer <your_access_token>` | JWT Access Token |

### Request Body

| Field | Type | Required | Description | Constraints |
| :--- | :--- | :--- | :--- | :--- |
| `workspace_id` | number | **Yes** | The ID of the workspace where the board will be created. | Must be a valid existing Workspace ID. |
| `name` | string | **Yes** | The title of the board. | |
| `images` | string | **Yes** | The background image URL or identifier. | |

**Example Request:**

```json
{
  "workspace_id": 1,
  "name": "Project Roadmap",
  "images": "https://images.unsplash.com/photo-123"
}
```

### Response

#### Success (200 OK)

Returns a success message upon creation.

**Example Response:**

```json
{
  "message": "Board created successfully"
  // Note: Currently does not return the created board object, only a success message.
}
```

#### Errors

| Status Code | Description | Example Body |
| :--- | :--- | :--- |
| **400 Bad Request** | Missing fields, invalid JSON, or invalid Workspace ID. | `{"error": "Key: 'Boards.Name' Error:Field validation for 'Name' failed..."}` |
| **401 Unauthorized** | Missing or invalid JWT token. | `{"error": "Unauthorized"}` |
| **500 Internal Server Error** | Server-side error (e.g., database failure). | `{"error": "Failed to create board..."}` |
