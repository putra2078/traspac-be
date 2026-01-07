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

## TaskCard Operations

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
Update content, name, date, or status of a card.

**Action**: `update_task_card`

**Payload**:
```json
{
  "task_card_id": 1,
  "name": "New Card Title",
  "content": "Updated content here",
  "date": "2023-12-31",
  "status": true
}
```
*Note: Omit fields that should not be updated.*

## TaskTab Operations

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


---

## Label Operations

### 8. Create Label
Add a new label to a task card.

**Action**: `create_label`

**Payload**:
```json
{
  "task_card_id": 1,
  "title": "Urgent",
  "color": "red"
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "create_label",
  "status": "success",
  "payload": {
    "task_card_id": 1,
    "title": "Urgent",
    "color": "red"
  },
  "data": {
    "id": 10,
    "task_card_id": 1,
    "title": "Urgent",
    "color": "red",
    "created_at": "2025-12-20T22:45:00Z",
    "updated_at": "2025-12-20T22:45:00Z"
  }
}
```

### 9. Get Labels
Fetch all labels for a specific task card.

**Action**: `get_labels`

**Payload**:
```json
{
  "task_card_id": 1
}
```

**Success Response** (sent only to requesting client):
```json
{
  "action": "get_labels",
  "status": "success",
  "payload": {
    "task_card_id": 1
  },
  "data": [
    {
      "id": 10,
      "task_card_id": 1,
      "title": "Urgent",
      "color": "red",
      "created_at": "2025-12-20T22:45:00Z",
      "updated_at": "2025-12-20T22:45:00Z"
    }
  ]
}
```

### 10. Update Label
Update an existing label's title or color.

**Action**: `update_label`

**Payload**:
```json
{
  "id": 10,
  "title": "Very Urgent",
  "color": "darkred"
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "update_label",
  "status": "success",
  "payload": {
    "id": 10,
    "title": "Very Urgent",
    "color": "darkred"
  },
  "data": {
    "id": 10,
    "task_card_id": 1,
    "title": "Very Urgent",
    "color": "darkred",
    "created_at": "2025-12-20T22:45:00Z",
    "updated_at": "2025-12-20T22:50:00Z"
  }
}
```

### 11. Delete Label
Delete a label.

**Action**: `delete_label`

**Payload**:
```json
{
  "id": 10
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "delete_label",
  "status": "success",
  "payload": {
    "id": 10
  },
  "data": {
    "id": 10
  }
}
```

---

## TaskCardUser Operations

### 12. Assign Task Card User
Assign a user to a task card.

**Action**: `assign_task_card_user`

**Payload**:
```json
{
  "task_card_id": 1,
  "user_id": 10
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "assign_task_card_user",
  "status": "success",
  "payload": {
    "task_card_id": 1,
    "user_id": 10
  },
  "data": {
    "id": 1,
    "task_card_id": 1,
    "user_id": 10,
    "created_at": "2025-12-23T20:25:00Z",
    "updated_at": "2025-12-23T20:25:00Z"
  }
}
```

### 13. Get Task Card Users
Fetch all users assigned to a specific task card.

**Action**: `get_task_card_users`

**Payload**:
```json
{
  "task_card_id": 1
}
```

**Success Response** (sent only to requesting client):
```json
{
  "action": "get_task_card_users",
  "status": "success",
  "payload": {
    "task_card_id": 1
  },
  "data": [
    {
      "id": 1,
      "task_card_id": 1,
      "user_id": 10,
      "created_at": "2025-12-23T20:25:00Z",
      "updated_at": "2025-12-23T20:25:00Z"
    }
  ]
}
```

### 14. Unassign Task Card User
Remove a user assignment from a task card.

**Action**: `unassign_task_card_user`

**Payload**:
```json
{
  "id": 1
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "unassign_task_card_user",
  "status": "success",
  "payload": {
    "id": 1
  },
  "data": {
    "id": 1
  }
}
```

---

## BoardUser Operations

### 15. Assign Board User
Assign a user to a board.

**Action**: `assign_board_user`

**Payload**:
```json
{
  "board_id": 1,
  "user_id": 10
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "assign_board_user",
  "status": "success",
  "payload": {
    "board_id": 1,
    "user_id": 10
  },
  "data": {
    "id": 1,
    "board_id": 1,
    "user_id": 10,
    "created_at": "2025-12-23T20:30:00Z",
    "updated_at": "2025-12-23T20:30:00Z"
  }
}
```

### 16. Get Board Users
Fetch all users assigned to a specific board.

**Action**: `get_board_users`

**Payload**:
```json
{
  "board_id": 1
}
```

**Success Response** (sent only to requesting client):
```json
{
  "action": "get_board_users",
  "status": "success",
  "payload": {
    "board_id": 1
  },
  "data": [
    {
      "id": 1,
      "board_id": 1,
      "user_id": 10,
      "created_at": "2025-12-23T20:30:00Z",
      "updated_at": "2025-12-23T20:30:00Z"
    }
  ]
}
```

### 17. Unassign Board User
Remove a user assignment from a board.

**Action**: `unassign_board_user`

**Payload**:
```json
{
  "id": 1
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "unassign_board_user",
  "status": "success",
  "payload": {
    "id": 1
  },
  "data": {
    "id": 1
  }
}
```

---

## WorkspaceUser Operations

### 18. Assign Workspace User
Assign a user to a workspace.

**Action**: `assign_workspace_user`

**Payload**:
```json
{
  "workspace_id": 1,
  "user_id": 10
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "assign_workspace_user",
  "status": "success",
  "payload": {
    "workspace_id": 1,
    "user_id": 10
  },
  "data": {
    "id": 1,
    "workspace_id": 1,
    "user_id": 10,
    "created_at": "2025-12-23T20:35:00Z",
    "updated_at": "2025-12-23T20:35:00Z"
  }
}
```

### 19. Get Workspace Users
Fetch all users assigned to a specific workspace.

**Action**: `get_workspace_users`

**Payload**:
```json
{
  "workspace_id": 1
}
```

**Success Response** (sent only to requesting client):
```json
{
  "action": "get_workspace_users",
  "status": "success",
  "payload": {
    "workspace_id": 1
  },
  "data": [
    {
      "id": 1,
      "workspace_id": 1,
      "user_id": 10,
      "created_at": "2025-12-23T20:35:00Z",
      "updated_at": "2025-12-23T20:35:00Z"
    }
  ]
}
```

### 20. Unassign Workspace User
Remove a user assignment from a workspace.

**Action**: `unassign_workspace_user`

**Payload**:
```json
{
  "id": 1
}
```

**Success Response** (broadcasted to all clients):
```json
{
  "action": "unassign_workspace_user",
  "status": "success",
  "payload": {
    "id": 1
  },
  "data": {
    "id": 1
  }
}
```
