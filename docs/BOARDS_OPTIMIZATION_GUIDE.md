# Boards API Optimization Guide

The Boards API has been optimized to improve performance and reduce memory usage. The following changes have been made:

## Key Changes
1. **Lightweight List**: `GET /api/v1/boards/` now returns only board metadata. It does **not** return Task Tabs or Task Cards.
2. **Lightweight Detail**: `GET /api/v1/boards/:id` now returns board metadata and Task Tabs. It does **not** return Task Cards.
3. **New Endpoints**: New endpoints have been added to fetch Task Tabs and Task Cards separately, with pagination for cards.

## API Usage

### 1. Get Board (Detail)
**Endpoint:** `GET /api/v1/boards/:id`
**Response:** Returns Board object with populated `task_tabs` but empty `task_cards`.

### 2. Get Tabs by Board
**Endpoint:** `GET /api/v1/boards/:board_id/tabs`
**Response:** Returns list of TaskTabs for the board.

### 3. Get Cards by Tab (Paginated)
**Endpoint:** `GET /api/v1/boards/tabs/:tab_id/cards`
**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50)
**Response:** Returns list of TaskCards for the tab, including Labels and Members snippets.

## Frontend Migration Guide
To adopt these changes, the frontend should:
1. Fetch the board detail.
2. Render the board and its tabs.
3. For each tab, trigger a separate (lazy/infinite scroll) request to `GET /api/v1/boards/tabs/:tab_id/cards` to load cards.
