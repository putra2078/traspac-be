# WebSocket API: Task Card Updates & Broadcasts

This document describes how to integrate real-time task card updates on the frontend. Use these actions to update cards and handle data streams when other users make changes.

---

## 1. WebSocket Actions

### 1.1 `update_task_card`
Use this action for general updates to a task card, including moving it to a different tab (useful for Drag & Drop).

**Payload Request:**
```json
{
  "action": "update_task_card",
  "payload": {
    "task_card_id": 101,
    "task_tab_id": 5,        // Optional: Move card to another tab
    "name": "Updated Title", // Optional
    "content": "New content", // Optional
    "date": "2024-12-31",    // Optional
    "status": true           // Optional: Completion status
  }
}
```

### 1.2 `update_task_tab_id`
Dedicated action for moving a task card to another tab.

**Payload Request:**
```json
{
  "action": "update_task_tab_id",
  "payload": {
    "task_card_id": 101,
    "task_tab_id": 5
  }
}
```

---

## 2. Broadcast Response (Success)

When an update is successful, the server broadcasts the **full updated task card** to all clients subscribed to the board.

**Broadcast Message:**
```json
{
  "action": "update_task_card", // or "update_task_tab_id"
  "status": "success",
  "payload": { ... original request payload ... },
  "data": {
    "id": 101,
    "task_tab_id": 5,
    "name": "Updated Title",
    "content": "...",
    "date": "...",
    "status": true,
    "labels": [...],    // Full preloaded labels
    "comments": [...],  // Full preloaded comments (with user)
    "members": [...],   // Full preloaded members (with user)
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

> [!IMPORTANT]
> **Data Consistency**: The `data` field contains the fresh state of the card after the update. Use this to replace the local card object in your state management store (e.g., Vuex, Pinia) to ensure the UI stays synchronized with all relations (comments, members, etc.).

---

## 3. Frontend Implementation Tips

### Handling Drag & Drop (Moving Tabs)
When a user finishes dragging a card from `Tab A` to `Tab B`:
1. Send `update_task_card` with the new `task_tab_id`.
2. Do **not** manually move the card in the UI if you want to rely on the server's confirmation.
3. Upon receiving the `update_task_card` broadcast, update your local state. If `task_tab_id` changed, the card should be re-rendered in the new column.

### Example Handler (JavaScript)
```javascript
function handleTaskCardUpdate(data) {
    const cardId = data.id;
    const newTabId = data.task_tab_id;
    
    // 1. Find the card and update its data
    updateCardInStore(cardId, data);
    
    // 2. If you use tab-based grouping, the UI will auto-update 
    // because data.task_tab_id now reflects the new position.
}
```

---

## 4. Verification Checklist
- [ ] Card name/status update reflects in real-time.
- [ ] Moving card between tabs updates all connected clients.
- [ ] Relations (Labels/Members) are preserved in the broadcast data.
