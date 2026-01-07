# WebSocket Repair & Integration Guide

## Overview
This document details all WebSocket payload fixes and provides complete frontend integration instructions to eliminate UI delays and ensure real-time updates work correctly.

## Problem Summary
**Root Cause**: WebSocket broadcasts were sending incomplete payloads without related `User` data, causing frontend rendering failures that appeared as "delays" until manual refresh.

**Solution**: All assignment/comment entities now include full `User` objects in broadcasts, enabling immediate UI updates.

---

## 1. Fixed Entities & Payloads

### 1.1 TaskCardUsers (Assign/Unassign Users to Cards)

**Actions**: `assign_task_card_user`, `unassign_task_card_user`

**New Payload Structure**:
```json
{
  "action": "assign_task_card_user",
  "status": "success",
  "payload": {
    "task_card_id": 101,
    "user_id": 5
  },
  "data": {
    "id": 12,
    "task_card_id": 101,
    "user_id": 5,
    "user": {
      "id": 5,
      "username": "johndoe",
      "email": "john@example.com"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 1.2 TaskCardComment (Create/Update/Delete Comments)

**Actions**: `create_task_card_comment`, `update_task_card_comment`, `delete_task_card_comment`

**New Payload Structure**:
```json
{
  "action": "create_task_card_comment",
  "status": "success",
  "payload": {
    "task_card_id": 101,
    "comment": "Great work!"
  },
  "data": {
    "id": 1,
    "task_card_id": 101,
    "user_id": 5,
    "user": {
      "id": 5,
      "username": "johndoe",
      "email": "john@example.com"
    },
    "comment": "Great work!",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 1.3 BoardsUsers (Assign/Unassign Users to Boards)

**Actions**: `assign_board_user`, `unassign_board_user`

**New Payload Structure**:
```json
{
  "action": "assign_board_user",
  "status": "success",
  "payload": {
    "board_id": 10,
    "user_id": 5
  },
  "data": {
    "id": 20,
    "board_id": 10,
    "user_id": 5,
    "user": {
      "id": 5,
      "username": "johndoe",
      "email": "john@example.com"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 1.4 WorkspacesUsers (Assign/Unassign Users to Workspaces)

**Actions**: `assign_workspace_user`, `unassign_workspace_user`

**New Payload Structure**: Same as BoardsUsers but with `workspace_id` instead of `board_id`.

---

## 2. Frontend Integration Instructions

### 2.1 WebSocket Message Handler Pattern

Update your WebSocket `onmessage` handler to process the `data` field:

```javascript
socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.status === 'success') {
    handleSuccessMessage(message);
  } else if (message.status === 'error') {
    handleErrorMessage(message);
  }
};

function handleSuccessMessage(msg) {
  switch (msg.action) {
    case 'assign_task_card_user':
      handleAssignTaskCardUser(msg.data);
      break;
    case 'create_task_card_comment':
      handleCreateComment(msg.data);
      break;
    case 'assign_board_user':
      handleAssignBoardUser(msg.data);
      break;
    // ... other actions
  }
}
```

### 2.2 Specific Handler Examples

#### A. TaskCardUsers Assignment
```javascript
function handleAssignTaskCardUser(data) {
  // data.user is now available!
  const assignment = {
    id: data.id,
    userId: data.user_id,
    userName: data.user.username,
    userAvatar: data.user.photo || '/default-avatar.png'
  };
  
  // Find the card and update its assignees
  const card = findCardById(data.task_card_id);
  if (card) {
    card.assignees.push(assignment);
    // UI will auto-update via reactivity
  }
}

function handleUnassignTaskCardUser(data) {
  const card = findCardById(data.task_card_id);
  if (card) {
    card.assignees = card.assignees.filter(a => a.id !== data.id);
  }
}
```

#### B. TaskCardComment Creation
```javascript
function handleCreateComment(data) {
  const comment = {
    id: data.id,
    text: data.comment,
    userId: data.user_id,
    userName: data.user.username,
    userAvatar: data.user.photo || '/default-avatar.png',
    createdAt: data.created_at
  };
  
  // Add to card's comments array
  const card = findCardById(data.task_card_id);
  if (card) {
    card.comments.unshift(comment); // Add to top
    // Reset loading state
    isSubmittingComment.value = false;
  }
}
```

#### C. BoardUsers Assignment
```javascript
function handleAssignBoardUser(data) {
  const member = {
    id: data.id,
    userId: data.user_id,
    userName: data.user.username,
    userEmail: data.user.email,
    userAvatar: data.user.photo || '/default-avatar.png'
  };
  
  // Add to board members list
  const board = findBoardById(data.board_id);
  if (board) {
    board.members.push(member);
  }
}
```

### 2.3 Critical: Reset Loading States

**IMPORTANT**: Always reset loading/submitting states in your handlers:

```javascript
// In your component
const isSubmittingComment = ref(false);

function createComment() {
  if (isSubmittingComment.value) return;
  
  isSubmittingComment.value = true;
  webSocketService.send('create_task_card_comment', {
    task_card_id: cardId,
    comment: commentText.value
  });
  commentText.value = '';
}

// In WebSocket handler
function handleCreateComment(data) {
  // ... add comment to UI ...
  
  // CRITICAL: Reset loading state
  isSubmittingComment.value = false;
}
```

### 2.4 Error Handling

Always handle WebSocket errors to reset UI states:

```javascript
function handleErrorMessage(msg) {
  console.error('WebSocket error:', msg.message);
  
  // Reset all loading states
  isSubmittingComment.value = false;
  isAssigningUser.value = false;
  
  // Show error to user
  showToast(msg.message, 'error');
}
```

---

## 3. Testing Checklist

After integration, verify:

- [ ] **TaskCardUsers**: Assign user → avatar appears immediately without refresh
- [ ] **TaskCardUsers**: Unassign user → avatar disappears immediately
- [ ] **Comments**: Create comment → appears immediately with user avatar
- [ ] **Comments**: Multiple rapid comments → all appear, no stuck loading state
- [ ] **BoardUsers**: Assign member → appears in members list immediately
- [ ] **WorkspaceUsers**: Assign member → appears immediately
- [ ] **Loading States**: All submit buttons re-enable after action completes
- [ ] **Error Cases**: Failed actions reset loading states and show error

---

## 4. Common Issues & Solutions

### Issue: "Stuck in loading state"
**Cause**: Loading state not reset in WebSocket handler  
**Solution**: Add `isSubmitting.value = false` in success/error handlers

### Issue: "User avatar not showing"
**Cause**: Frontend trying to access `data.user` before fix  
**Solution**: Update to use `msg.data.user` as shown in examples above

### Issue: "Changes only appear after refresh"
**Cause**: Not listening to WebSocket events or not updating reactive state  
**Solution**: Ensure WebSocket handlers update the same reactive objects used by UI

---

## 5. Migration Steps

1. **Update WebSocket Service**: Modify message handler to use `msg.data` instead of `msg.payload`
2. **Update Action Handlers**: Add handlers for all actions listed in Section 1
3. **Add User Data**: Update UI components to display `data.user.username`, `data.user.photo`, etc.
4. **Reset States**: Add loading state resets in all handlers
5. **Test Thoroughly**: Follow testing checklist in Section 3

---

## Support

If issues persist after following this guide:
1. Check browser console for WebSocket errors
2. Verify `msg.data.user` exists in received messages
3. Confirm reactive state updates are triggering UI re-renders
4. Check that `join_board` action is sent when opening a board
