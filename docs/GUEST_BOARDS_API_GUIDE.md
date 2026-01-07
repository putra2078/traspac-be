# Guest Boards Fetching Guide

This guide provides information on how to fetch Boards where the current user is a member but not the creator (Guest Boards) using the REST API.

## API Endpoint

To retrieve all boards associated with the authenticated user (both owned and guest boards), use the following endpoint:

- **URL**: `/api/v1/boards/`
- **Method**: `GET`
- **Auth Required**: YES (JWT Bearer Token)

## Identification of Guest Boards

Since the API returns both owned boards and guest boards in a single list, the frontend should distinguish them by comparing the `created_by` field with the current user's ID.

| Board Type | Logic |
| :--- | :--- |
| **Owned Board** | `board.created_by === current_user_id` |
| **Guest Board** | `board.created_by !== current_user_id` |

## Request Example

### Axios (Frontend)
```javascript
const response = await axios.get('/api/v1/boards/', {
  headers: {
    Authorization: `Bearer ${localStorage.getItem('token')}`
  }
});

const allBoards = response.data.data;
const guestBoards = allBoards.filter(board => board.created_by !== currentUserId);
const ownedBoards = allBoards.filter(board => board.created_by === currentUserId);
```

## Response Structure

The response will be a JSON object containing an array of boards. Each board entry includes task tabs and task cards preloaded with labels and members.

```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "workspace_id": 1,
      "created_by": 2, 
      "name": "Project Alpha",
      "images": "banner.jpg",
      "task_tabs": [...],
      "task_cards": [
        {
          "id": 101,
          "name": "Implement Login",
          "labels": [...],
          "members": [...]
        }
      ],
      "created_at": "2025-12-30T10:00:00Z",
      "updated_at": "2025-12-30T10:00:00Z"
    }
  ]
}
```

## Backend Implementation Note

The backend implementation for this logic resides in:
- **Handler**: `internal/domain/boards/handler.go` -> `GetByUserID`
- **UseCase**: `internal/domain/boards/usecase.go` -> `FindByUserID`

The `FindByUserID` method performs the following:
1. Fetches boards where `created_by` matches the `userID`.
2. Fetches board associations from the `boards_users` table for the `userID`.
3. Merges and deduplicates the results.
4. Enriches each board with its associated `TaskTabs` and `TaskCards` (including bulk-fetched `Labels` and `Members`).
