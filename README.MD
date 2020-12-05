# ChitChat
ChitChat is a simple chat application that helps you connect with other people around the world.

## Install ##
All that you need is [Golang](https://golang.org/). Once you run the application, it will expose [port 8080](./config.json) on the host.
```
20:5:13 app         | ChitChat 0.1 started at 127.0.0.1:8080
```

## Configuration ##
Edit [config.json](./config.json) to configure HTTP server settings.

## ChatRoom API ##
Send different HTTP requests to `/chat/`:
  * `GET /chats/<room_id>`: retrieve a chat room by ID or title
  * `POST /chats`: create a new chat room
  * `PUT /chats/<room_id>`: update a chat room by ID or title
  * `DEL /chats/<room_id>`: delete a chat room

E.g. A basic chat room POST request could look as follows:
```json
{
   "title":"My chat room",
   "description":"There are many like it, but this one is mine",
   "visibility":"public"
}
```
If successful, server should respond with code 201 and the newly created resource:
```json
{
    "title": "My chat room",
    "description": "There are many like it, but this one is mine",
    "visibility": "public",
    "createdAt": "2020-12-03T21:23:54.4213184-08:00",
    "updatedAt": "2020-12-03T21:23:54.4213184-08:00",
    "id": 2,
    "users": []
}
```

If unsuccessful, the server will return an error, e.g. if the room already exists: 
```json
{
    "status": false,
    "error": {
        "code": 102,
        "error": "Room error: Duplicate room",
        "field": "My chat room"
    }
}
```
The minimum required fields to create a public chat room are `"title"` and `"visibility"`. "Private" and "hidden" chat rooms also require a password.

## Chatting ##
Chatting is implemented using Server-Sent-Events (SSE) as a step-by-step procedure. Any endpoints that mandates `text/event-stream` content type as opposed to JSON are behind the `chats/<room_id>/sse/` subpath. The steps are:
  1. Send a HTTP POST with room password to `/chats/<room_id>/token` to obtain authorization (only required for private and hidden rooms)
  2. Subscribe to chat event stream via HTTP GET request to `/chats/<room_id>/sse/subscribe`
  3. Send a HTTP POST to `/chats/<room_id>/sse/broadcast` to broadcast messages using the following format:
  ```
  {
      "event_type": "join/leave/send",
      "name": "my-username",
      "color": "color of chat message",
      "msg": "message to send",
      "secret": "password" # (only needed for private or secret rooms token request)
  }
  ```
  It is recommended to send a `join` and `leave` event broadcast when joining/leaving chat rooms so everyone else is notified.
  4. Send a HTTP GET to `/chats/<room_id>/token/renew` to renew the authorization token before expiration.

  ## TODO ##
  1. Need to add more unit tests for chatting. 
    * Also need to add mocks once DB is setup for decoupled testing
  2. Passwords must be salted and authorization improved wit expiration and retries
    * I also plan to use JWT instead of session-based cookie approach
  3. Generate proper UUIDs instead of integers as IDs
  4. Need to move chat server data layer into DB, right now everything is stored in memory to allow for rapid development and testing
  5. Need to improve reliability and longevity of chatting sessions, on some browsers (Firefox) there are intermittent connection retries.

  ## Auth:
  ### Session-based:
  Client-side:
  * hash password before sending it (optional)
  * Loop to refresh once per minute in chat
  Server-side:
  * Refactor to create auth handler
  * Add session ID to client
  * Add session refresh handler
  * Hide password, session ID from any public-facing API response