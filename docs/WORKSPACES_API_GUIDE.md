# Workspaces API Guide

This document outlines the REST API endpoints for managing workspaces.

## 1. Create Workspace
Create a new workspace.
- **Endpoint:** `POST /api/v1/workspaces/`
- **Request Body:**
```json
{
  "name": "My Workspace",
  "privacy": "public" // options: public, private, team
}
```
- **Response:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "My Workspace",
    "created_by": 5,
    "pass_code": "ASDFGH",
    "join_link": "http://be.putratek.my.id/api/v1/workspaces/join?token=..."
  }
}
```

## 2. Get My Workspaces
Fetch workspaces created by or where the user is a member (Full Access).
- **Endpoint:** `GET /api/v1/workspaces/`
- **Response:** List of workspaces.

## 3. Get Guest Workspaces
Fetch workspaces where the user is a member BUT not the owner (Guest Access).
- **Endpoint:** `GET /api/v1/workspaces/guest`
- **Response:** List of guest workspaces.

## 4. Get Workspace by ID
- **Endpoint:** `GET /api/v1/workspaces/:id`
- **Response:** Workspace details.

## 5. Join Workspace (via Token)
- **Endpoint:** `POST /api/v1/workspaces/join?token={JWT_TOKEN}`
- **Note:** This is usually triggered by the Join Link.

## 6. Generate Join Token
- **Endpoint:** `GET /api/v1/workspaces/:id/join-token`
- **Response:** `{ "token": "..." }`
