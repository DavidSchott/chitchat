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
  * `GET /chat/<room_id>`: retrieve a chat room by ID or title
  * `POST /chat/`: create a new chat room
  * `PUT /chat/<room_id>`: update a chat room by ID or title
  * `DEL /chat/<room_id>`: delete a chat room

E.g. a sample chat room request body could look as follows:
```
{
   "title":"My chat room",
   "description":"There are many like it, but this one is mine",
   "visibility":"public"
}
```
The minimum required fields to create a public chat room are title and visibility. "Private" and "hidden" chat rooms also require a password.

## Chatting ##
Chatting is implemented using Server-Sent-Events (SSE) as a step-by-step procedure:
  1. Send a HTTP POST to `/chat/sse/login` to start a session (only needed for private and hidden rooms)
  2. Subscribe to chat notifications by listening to `/chat/sse` source events
  3. Send a HTTP POST to `/chat/sse/event` to broadcast messages using the following JSON format:
  ```
  {
      type: "join/leave/send",
      name: "my-username",
      id: "room_id",
      color: "color of chat message",
      msg: "message to send",
      secret: "password" # (only needed for private or secret rooms)
  }
  ```

  ## TODO ##
  1. Need to add more unit tests for chatting. 
    * Also need to add mocks once DB is setup for decoupled testing
  2. Passwords must be salted on client
  3. Generate proper UUIDs instead of integers as IDs
  4. Need to move chat server data layer into DB, right now everything is stored in memory to allow for rapid development and testing
  5. Need to improve reliability and longevity of chatting sessions, on some browsers (Firefox) there are random interruptions