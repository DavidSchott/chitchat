# ChitChat
![Go](https://github.com/DavidSchott/chitchat/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/DavidSchott/chitchat)](https://goreportcard.com/report/github.com/DavidSchott/chitchat)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/DavidSchott/chitchat/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/DavidSchott/chitchat/branch/master/graph/badge.svg?token=G6AOS5QVHR)](https://codecov.io/gh/DavidSchott/chitchat)

ChitChat is a simple chat app with basic features such as password protection and real-time communication.

## Install ##
All that you need is [Golang](https://golang.org/). Once you run the application, it will expose a [target port](./config.json) on the host.
```
20:5:13 app         | ChitChat 0.3 started at 127.0.0.1:443
```

## Configuration ##
Edit [config.json](./config.json) to configure HTTP server settings.

## ChatRoom API ##
Send HTTP requests to `/chats`:
  * `GET /chats/<room_id>`: retrieve a chat room by ID or title
  * `POST /chats`: create a new chat room
  * `PUT /chats/<room_id>`: update existing chat room by ID or title
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

If unsuccessful, the server will return an error with a [custom code](./data/apierror.go), e.g. if a room with the same title already exists: 
```json
{
    "status": false,
    "error": {
        "code": 102,
        "error": "Room error: Duplicate room",
        "field": "title"
    }
}
```
The minimum required fields to create a public chat room are `"title"` and `"visibility"`. "Private" and "hidden" chat rooms also require a password.

  ## Authorization
  Authorization is implemented using [JSON Web Tokens](https://jwt.io/introduction/) in an attempt to achieve statelessness and scalability. In order to authorize chatting access to a password-protected chat room, users will need to pass in a token in the `Authorization` HTTP header using the Bearer scheme. 
  
  The token can be requested by sending a JSON body containing the password to the `/chats/<room_id>/token` endpoint. For example:
  ```json
  {
   "secret":"my_secret",
   "name":"david"
  }
  ```
  If successful, the server will respond with the token in the response body.
```json
{
    "status": true,
    "name": "david",
    "room_id": 2,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImRhdmlkIiwicm9vbV9pZCI6MiwiZXhwIjoxNjA3NTg1MjQ5fQ.b6XnNqrFnFmuUMhTBKfyR3PAyCQkxbUaPupBXgknl8w"
}
```  
Otherwise an error will be returned if the supplied password is invalid or missing.
```json
{
    "status": false,
    "error": {
        "code": 304,
        "error": "Unauthorized operation",
        "field": "secret"
    }
}
```
For added convenience, it is also possible to request a new token before the current one expires by sending an authorized request to `/chats/<room_id>/token/renew`.

## Chatting ##
Chatting is implemented using server-sent events (SSE) as a step-by-step procedure. Any endpoints that mandate `text/event-stream` content type as opposed to JSON are behind the `chats/<room_id>/sse/` subpath. The steps are:
  1. Send a HTTP POST to `/chats/<room_id>/token` to obtain [authorization](#authorization) (only required for private and hidden rooms)
  2. Subscribe to chat event stream via HTTP GET request to `/chats/<room_id>/sse/subscribe`
  3. Send a HTTP POST to `/chats/<room_id>/sse/broadcast` to broadcast messages in the following format:
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
  1. Need to add more unit tests for chatting 
      * Also need to add test mocks once DB is setup for decoupled testing from model/data layer. Will likely inject dependencies in top-down fashion.
  2. Authorization token handling on client-side website frontend could be improved with retries
      * Send requests to refresh token shortly before expiration (e.g. n-5 minutes? What's a good tolerance?)
      * Hide password field showing the salted hash from public-facing API response?
  3. Generate proper UUIDs instead of integers as IDs?
  4. Need to store chat server data models in a DB, right now everything is stored in memory to allow for rapid development and testing
  5. Need to improve reliability of chatting sessions
      * Adding heartbeat pings as keep-alives
      * Can also move to WebSockets instead of SSE
  6. Switch to a proper logging library such as logrus