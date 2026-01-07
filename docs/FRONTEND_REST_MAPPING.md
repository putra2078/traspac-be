# Frontend Refactor Guide: Transitioning Initial Load to REST API

This document provides a mapping for the frontend to transition from WebSocket-based data queries to REST API endpoints. All "Initial Load" logic has been removed from the WebSocket layer to improve performance and architecture.

## WebSocket Actions Removed
The following WebSocket actions are **no longer supported**. If you send them, the server will ignore them or return an "Unknown action" error.

1. `get_board_users`
2. `get_task_card_users`
3. `get_task_card_comments`
4. `get_labels`
5. `get_workspace_users`

---

## REST API Mapping (Consolidated)

The `GET /api/v1/task-cards/:id` endpoint has been optimized to return all related data in a single request. 

> [!TIP]
> **Consolidated Initial Load**: You now only need to call `GET /api/v1/task-cards/:id` once to get the card, its labels, comments (with user details), and members (with user details).

| Entity | Field in TaskCard JSON | Nested Data Included |
| :--- | :--- | :--- |
| `Labels` | `labels` | Color, Title |
| `Comments` | `comments` | Comment text, User object |
| `Members` | `members` | User object |

### Mapping Table

| Data Type | Corresponding REST Endpoint (GET) | Notes |
| :--- | :--- | :--- |
| **Full TaskCard** | `/api/v1/task-cards/:id` | **Recommended** (Includes labels, comments, members) |
| `labels` | `/api/v1/labels/task-card/:id` | Still available if needed separately |
| `comments` | `/api/v1/task-card-comments/task-card/:id` | Still available if needed separately |
| `members` | `/api/v1/task-card-users/task-card/:id` | Still available if needed separately |

---

## Implementation Notes

### 1. Authentication
All REST endpoints above require the `Authorization: Bearer <token>` header (standard JWT authentication).

### 2. WebSocket Usage
WebSockets should now **only** be used for:
- Joining a board/room (`join_board`, `join_room_chat`).
- Real-time events (broadcasts of `create`, `update`, `delete`, `assign`, etc.).
- Sending real-time messages.

### 3. Flow Example
**Old Flow:**
1. Connect WS.
2. Send `get_labels` via WS.
3. Wait for WS response.

**New Recommended Flow:**
1. Call `GET /api/v1/labels/task-card/:id` via HTTP.
2. Populate the UI with the fetched data.
3. Rely on WS broadcasts to update the UI when other users make changes.
