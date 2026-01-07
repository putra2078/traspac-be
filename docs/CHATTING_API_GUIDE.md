# Chatting Feature API Guide (Workspace-based)

This document explains how to integrate the new workspace-based chatting feature, which uses a combination of REST API for initial data and WebSocket/Kafka for real-time messaging.

---

## 1. REST API (Initial Load)

Use these endpoints to fetch the list of rooms and memberships when a user enters a workspace.

### 1.1 Create Chat Room
Use this endpoint to create a new chat room within a workspace.

**Endpoint:** `POST /api/v1/room-chats/`
- **Request Body:**
```json
{
  "workspace_id": 1,
  "name": "General Discussion",
  "passcode": "123456" // Optional
}
```
- **Response:**
```json
{
  "status": "success",
  "data": {
    "id": 123,
    "workspace_id": 1,
    "name": "General Discussion",
    "created_by": 5,
    "link_join": "...",
    "created_at": "..."
  }
}
```

### 1.2 Fetch Rooms in Workspace
**Endpoint:** `GET /api/v1/room-chats/workspace/:workspace_id`
- **Response:** List of rooms (with `id`, `name`, `workspace_id`, etc.)

### 1.2 Fetch Room Members
**Endpoint:** `GET /api/v1/room-users/room/:room_id`
- **Response:** List of users joined to the room.

---

## 2. WebSocket & Real-time Messaging

The WebSocket connection handles room joining (automatic membership) and real-time message broadcasting via Kafka.

### 2.1 Joining a Room (`join_room_chat`)
Send this action to start receiving messages for a specific room.
**Request Payload:**
```json
{
  "action": "join_room_chat",
  "payload": {
    "room_id": 123
  }
}
```
**Behavior:**
- The server will automatically add the user to the `room_users` table if they aren't already a member.
- The server returns the **Chat History** (last 50 messages) in the success response.

### 2.2 Sending a Message (`send_room_chat_message`)
**Request Payload:**
```json
{
  "action": "send_room_chat_message",
  "payload": {
    "room_id": 123,
    "message_text": "Hello world!",
    "message_content": "" // Optional: for attachments or rich content
  }
}
```
**Behavior:**
- Message is saved to the database.
- Message is **immediately** broadcasted locally to current instance clients (no delay).
- Message is produced to Kafka to reach clients on other instances.
- Kafka consumer handles cross-instance broadcasting.

### 2.3 Typing Indicator (`typing_indicator`)
Send this action when the user starts or stops typing.
**Request Payload:**
```json
{
  "action": "typing_indicator",
  "payload": {
    "room_id": 123,
    "is_typing": true
  }
}
```
**Response (Broadcast to room):**
```json
{
  "action": "typing_indicator",
  "status": "success",
  "data": {
    "room_id": 123,
    "user_id": 5,
    "is_typing": true
  }
}
```

### 2.4 Receiving a Message (`new_room_chat_message`)
Incoming broadcast from server:
```json
{
  "action": "new_room_chat_message",
  "status": "success",
  "data": {
    "id": 456,
    "user_id": 1,
    "user": { ... user details ... },
    "room_id": 123,
    "message_text": "Hello world!",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 3. Implementation Workflow
1. **Initial Load**: Call `GET /api/v1/room-chats/workspace/:id` to show the room list.
2. **Select Room**: When a user clicks a room, call `join_room_chat` via WebSocket.
3. **Display History**: Render the messages received in the `join_room_chat` success response.
4. **Listen**: Append incoming messages with action `new_room_chat_message` to the UI state.
