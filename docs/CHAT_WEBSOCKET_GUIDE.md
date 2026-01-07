# Chat WebSocket Integration Guide

This document describes how to integrate the enhanced Chat WebSocket features into the frontend. The chat system now provides user identity information (names) for both messages and typing indicators.

## Connection

Connect to the WebSocket endpoint using the JWT token:

```javascript
const socket = new WebSocket('ws://be.putratek.my.id/api/v1/ws/task-cards?token=YOUR_JWT_TOKEN');
```

## Chat Room Actions

### 1. Join Chat Room
**Action:** `join_room_chat`

**Request:**
```json
{
  "action": "join_room_chat",
  "payload": {
    "room_id": 1
  }
}
```

**Response (History):**
The server will return the chat history. Each message now includes a `user` object with the sender's details.
```json
{
  "action": "join_room_chat",
  "status": "success",
  "data": [
    {
      "id": 10,
      "user_id": 1,
      "user": {
        "id": 1,
        "username": "johndoe",
        "email": "john@example.com"
      },
      "room_id": 1,
      "message_text": "Hello everyone!",
      "created_at": "2026-01-01T10:00:00Z"
    }
  ]
}
```

### 2. Send Message
**Action:** `send_room_chat_message`

**Request:**
```json
{
  "action": "send_room_chat_message",
  "payload": {
    "room_id": 1,
    "message_text": "Hello!",
    "message_content": "" 
  }
}
```

**Broadcast (New Message):**
All members in the room (including the sender) will receive this broadcast.
```json
{
  "action": "new_room_chat_message",
  "status": "success",
  "data": {
    "id": 11,
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com"
    },
    "room_id": 1,
    "message_text": "Hello!",
    "created_at": "2026-01-01T10:05:00Z"
  },
  "sender_name": "John Doe",
  "sender_username": "johndoe"
}
```

### 3. Typing Indicator
**Action:** `typing_indicator`

**Request:**
Send this when the user starts or stops typing.
```json
{
  "action": "typing_indicator",
  "payload": {
    "room_id": 1,
    "is_typing": true
  }
}
```

**Broadcast:**
Other users in the room will receive the typing status along with the user's name.
```json
{
  "action": "typing_indicator",
  "status": "success",
  "data": {
    "room_id": 1,
    "user_id": 1,
    "user_name": "John Doe",
    "is_typing": true
  }
}
```

### 4. Upload File (Attachment)
**Endpoint:** `POST /api/v1/room-chats/upload`
**Header:** `Authorization: Bearer <token>`
**Content-Type:** `multipart/form-data`

**Request:**
Form Data:
- `file`: (binary file)

**Response:**
```json
{
  "status": "success",
  "data": {
    "url": "https://<supabase-project>.supabase.co/storage/v1/object/public/chat-attachments/room-chats/uuid.jpg"
  }
}
```

**Usage Flow:**
1.  Frontend uploads the file to `POST /api/v1/room-chats/upload`.
2.  Server returns the public URL of the uploaded file.
3.  Frontend sends the message via WebSocket using the `send_room_chat_message` action, putting the returned URL into the `message_content` field.

```json
{
  "action": "send_room_chat_message",
  "payload": {
    "room_id": 1,
    "message_text": "Check this photo",
    "message_content": "https://<supabase-project>.supabase.co/storage/v1/object/public/chat-attachments/room-chats/uuid.jpg"
  }
}
```

### 5. Edit Message
**Action:** `edit_room_chat_message`

**Request:**
```json
{
  "action": "edit_room_chat_message",
  "payload": {
    "room_id": 1,
    "message_id": 15,
    "message_text": "Updated message text"
  }
}
```

**Broadcast:**
```json
{
  "action": "edit_room_chat_message",
  "status": "success",
  "data": {
    "id": 15,
    "room_id": 1,
    "message_text": "Updated message text",
    ...
  },
  "sender_name": "...",
  "sender_username": "..."
}
```

### 6. Delete Message
**Action:** `delete_room_chat_message`

**Request:**
```json
{
  "action": "delete_room_chat_message",
  "payload": {
    "room_id": 1,
    "message_id": 15
  }
}
```

**Broadcast:**
```json
{
  "action": "delete_room_chat_message",
  "status": "success",
  "data": {
    "message_id": 15,
    "room_id": 1
  }
}
```

## Summary of Changes
- **New Message Broadcast:** Now includes a `user` object preloaded from the database.
- **Typing Indicator Broadcast:** Now includes a `user_name` field (fetched from the contact/profile).
- **History:** The `join_room_chat` history now also preloads the `user` object for each message.
- **File Upload:** New REST endpoint `POST /api/v1/room-chats/upload` to handle file uploads. Files are stored in the configured Supabase Storage bucket (default: `user-uploads` or as configured in `config.yaml`) under the `room-chats` folder. Use the returned URL in `message_content`.
- **Edit/Delete:** Added `edit_room_chat_message` and `delete_room_chat_message` actions.
