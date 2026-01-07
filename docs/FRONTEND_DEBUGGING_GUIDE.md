# Frontend Debugging Guide: Resolving Fetching Loops & Memory Leaks

This guide helps identifying if your frontend execution environment (browser client) is causing server-side memory leaks due to infinite fetching loops or aggressive retry logic.

## Symptoms
- Server RAM increases even after user stops interacting.
- Logs show repeated `Rate limit exceeded` or `Context deadline exceeded`.
- Browser console shows rapid stream of network requests.

## 1. Check for Infinite Fetch Loops in React/Vue/JS

### React `useEffect` Pitfall
A common bug is adding the *fetched data* itself as a dependency to the effect that fetches it, causing an infinite loop.

**Bad Code:**
```javascript
useEffect(() => {
  fetchMessages();
}, [messages]); // ERROR: Fetching updates 'messages', which triggers fetch again!
```

**Fix:**
Remove the data from dependency array, or use a specific trigger (ID, status).
```javascript
useEffect(() => {
  fetchMessages();
}, [roomId]); // OK: Only fetches when room changes
```

### Vue `watch` Pitfall
Watching a complex object and modifying it inside the watcher.

**Bad Code:**
```javascript
watch(messages, () => {
  refreshMessages(); // Infinite loop if refreshMessages mutates 'messages'
}, { deep: true });
```

## 2. Check WebSocket Retry Logic

If your WebSocket connection drops (e.g., due to Rate Limit or Server Restart), does your frontend retry immediately without backing off?

**Bad Logic:**
```javascript
socket.onclose = () => {
    connect(); // Instant retry -> Spam -> Server bans -> Instant retry...
}
```

**Good Logic (Exponential Backoff):**
```javascript
let retries = 0;
socket.onclose = () => {
    const delay = Math.min(1000 * 2 ** retries, 30000); // 1s, 2s, 4s... max 30s
    setTimeout(() => {
        connect();
        retries++;
    }, delay);
}
```

## 3. Verify in Browser DevTools

1.  Open **Developer Tools (F12)** -> **Network** tab.
2.  Filter by **WS** (WebSockets) or **Fetch/XHR**.
3.  Look for:
    -   Requests appearing multiple times per second.
    -   A WebSocket connection that constantly connects and disconnects (Status 101 then Red).
4.  **Console** tab: Look for repetitive error logs.

## 4. Specific to This Issue (Chat Delay & RAM)

If you see `Error producing to Kafka: context deadline exceeded` on the server *after* you stop typing, it means the server has a **backlog** of messages to process.

*   **Cause:** Frontend sent too many messages in a burst (or loop).
*   **Fix:** Ensure your frontend **Debounces** input events (don't send on every keystroke if doing "typing indicator" logic without throttling).

**Example Throttle (Typing Indicator):**
```javascript
function sendTyping() {
    socket.send("typing");
}
const throttledTyping = _.throttle(sendTyping, 2000); // Only send once every 2s
```
