# WebSocket Usage Guide

This guide explains how to integrate the WebSocket API for real-time updates on `TaskCards` and `TaskTabs`.

## Endpoint
`ws://localhost:8080/api/v1/ws/task-cards`

## Connection
Connect to the detailed endpoint using a standard WebSocket client.

## Message Structure
All messages sent to the server must follow this JSON structure:

```json
{
  "action": "ACTION_NAME",
  "payload": { ... }
}
```

### 1. Update Task Card Tab
Move a card to a different tab.

**Action**: `update_task_tab_id`

**Payload**:
```json
{
  "task_card_id": 1,
  "task_tab_id": 5
}
```

### 2. Update Task Card Details
Update content, name, comment, date, or status of a card.

**Action**: `update_task_card`

**Payload**:
```json
{
  "task_card_id": 1,
  "name": "New Card Title",
  "content": "Updated content here",
  // "comment": "New comment",
  "date": "2023-12-31",
  "status": true
}
```
*Note: Omit fields that should not be updated.*

### 3. Update Task Tab Details
Update name or position of a tab.

**Action**: `update_task_tab`

**Payload**:
```json
{
  "task_tab_id": 2,
  "name": "Review",
  "position": 3
}
```

## Responses
The server broadcasts a success message to all connected clients when an update occurs.

**Success Response**:
```json
{
  "action": "update_task_card",
  "status": "success",
  "payload": { ... },
  "data": { ...updated_object_data... }
}
```

## Errors
If an error occurs, the server sends an error message to the specific client.

**Error Response**:
```json
{
  "status": "error",
  "error": "Error description"
}
```

---

## TaskCardComment Operations

### 4. Create Task Card Comment
Add a new comment to a task card.

**Action**: `create_task_card_comment`

**Payload**:
```json
{
  "task_card_id": 1,
  "comment": "This is a new comment"
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "create_task_card_comment",
  "status": "success",
  "payload": {
    "task_card_id": 1,
    "comment": "This is a new comment"
  },
  "data": {
    "id": 5,
    "task_card_id": 1,
    "comment": "This is a new comment",
    "created_at": "2025-12-19T23:45:00Z",
    "updated_at": "2025-12-19T23:45:00Z"
  }
}
```

### 5. Get Task Card Comments
Fetch all comments for a specific task card.

**Action**: `get_task_card_comments`

**Payload**:
```json
{
  "task_card_id": 1
}
```

**Success Response** (sent only to requesting client):
```json
{
  "action": "get_task_card_comments",
  "status": "success",
  "payload": {
    "task_card_id": 1
  },
  "data": [
    {
      "id": 1,
      "task_card_id": 1,
      "comment": "First comment",
      "created_at": "2025-12-19T10:00:00Z",
      "updated_at": "2025-12-19T10:00:00Z"
    },
    {
      "id": 2,
      "task_card_id": 1,
      "comment": "Second comment",
      "created_at": "2025-12-19T11:00:00Z",
      "updated_at": "2025-12-19T11:00:00Z"
    }
  ]
}
```

### 6. Update Task Card Comment
Update an existing comment.

**Action**: `update_task_card_comment`

**Payload**:
```json
{
  "id": 5,
  "comment": "Updated comment text"
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "update_task_card_comment",
  "status": "success",
  "payload": {
    "id": 5,
    "comment": "Updated comment text"
  },
  "data": {
    "id": 5,
    "task_card_id": 1,
    "comment": "Updated comment text",
    "created_at": "2025-12-19T23:45:00Z",
    "updated_at": "2025-12-19T23:50:00Z"
  }
}
```

### 7. Delete Task Card Comment
Delete a comment.

**Action**: `delete_task_card_comment`

**Payload**:
```json
{
  "id": 5
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "delete_task_card_comment",
  "status": "success",
  "payload": {
    "id": 5
  },
  "data": {
    "id": 5
  }
}
```

