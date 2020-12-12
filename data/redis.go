package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Channel name to use with redis
	Channel = "chat"
)

var (
	WaitingMessage, AvailableMessage []byte
	waitSleep                        = time.Duration(10 * int64(time.Minute))
)

func init() {
	var err error
	WaitingMessage, err = json.Marshal(ChatEvent{
		User:  "system",
		Msg:   "Waiting for redis to be available. Messaging won't work until redis is available",
		Color: "purple",
	})
	if err != nil {
		panic(err)
	}
	AvailableMessage, err = json.Marshal(ChatEvent{
		User:  "system",
		Msg:   "Redis is now available & messaging is now possible",
		Color: "purple",
	})
	if err != nil {
		panic(err)
	}
}

// RedisReceiver receives events from Redis and broadcasts them to all
// registered websocket connections that are Registered.
type RedisReceiver struct {
	pool *redis.Pool

	events         chan []byte
	newConnections chan *websocket.Conn
	rmConnections  chan *websocket.Conn
}

// newRedisReceiver creates a RedisReceiver that will use the provided
// rredis.Pool.
func NewRedisReceiver(pool *redis.Pool) RedisReceiver {
	return RedisReceiver{
		pool:           pool,
		events:         make(chan []byte, 1000), // TODO: 1000 is arbitrary, determine better threshold?
		newConnections: make(chan *websocket.Conn),
		rmConnections:  make(chan *websocket.Conn),
	}
}

func (rr *RedisReceiver) Wait(_ time.Time) error {
	rr.Broadcast(WaitingMessage)
	time.Sleep(waitSleep)
	return nil
}

// Run receives pubsub events from Redis after establishing a connection.
// When a valid message is received it is broadcast to all connected websockets
func (rr *RedisReceiver) Run() error {
	fmt.Println("Channel", Channel)
	conn := rr.pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(Channel)
	go rr.connHandler()
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Println("Redis Message Received:", v.Data)
			if _, err := validateEvent(v.Data); err != nil {
				fmt.Println("Error unmarshalling message from Redis", err.Error())
				continue
			}
			rr.Broadcast(v.Data)
		case redis.Subscription:
			fmt.Println("Redis subscription received", v.Kind, v.Count)
		case error:
			return errors.Wrap(v, "Error while subscribed to Redis channel")
		default:
			fmt.Println("Unknown Redis receive during subscription", v)
		}
	}
}

// Broadcast the provided message to all connected websocket connections.
// If an error occurs while writting a message to a websocket connection it is
// closed and deregistered.
func (rr *RedisReceiver) Broadcast(evt []byte) {
	rr.events <- evt
}

// register the websocket connection with the receiver.
func (rr *RedisReceiver) register(conn *websocket.Conn) {
	rr.newConnections <- conn
}

// deRegister the connection by closing it and removing it from our list.
func (rr *RedisReceiver) deRegister(conn *websocket.Conn) {
	rr.rmConnections <- conn
}

func (rr *RedisReceiver) connHandler() {
	conns := make([]*websocket.Conn, 0)
	for {
		select {
		case evt := <-rr.events:
			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, evt); err != nil {
					fmt.Println("Error writing data to connection. Closing and removing connection", evt, err, conn)
					conns = removeConn(conns, conn)
				}
			}
		case conn := <-rr.newConnections:
			conns = append(conns, conn)
		case conn := <-rr.rmConnections:
			conns = removeConn(conns, conn)
		}
	}
}

func removeConn(conns []*websocket.Conn, remove *websocket.Conn) []*websocket.Conn {
	var i int
	var found bool
	for i = 0; i < len(conns); i++ {
		if conns[i] == remove {
			found = true
			break
		}
	}
	if !found {
		fmt.Println(fmt.Sprintf("conns: %#v\nconn: %#v\n", conns, remove))
		panic("Conn not found")
	}
	copy(conns[i:], conns[i+1:]) // shift down
	conns[len(conns)-1] = nil    // nil last element
	return conns[:len(conns)-1]  // truncate slice
}

// RedisWriter publishes events to the Redis CHANNEL
type RedisWriter struct {
	pool   *redis.Pool
	events chan []byte
}

func NewRedisWriter(pool *redis.Pool) RedisWriter {
	return RedisWriter{
		pool:   pool,
		events: make(chan []byte, 10000),
	}
}

// Run the main RedisWriter loop that publishes incoming events to Redis.
func (rw *RedisWriter) Run() error {
	conn := rw.pool.Get()
	defer conn.Close()

	for data := range rw.events {
		if err := writeToRedis(conn, data); err != nil {
			rw.publish(data) // attempt to redeliver later
			return err
		}
	}
	return nil
}

func writeToRedis(conn redis.Conn, data []byte) error {
	if err := conn.Send("PUBLISH", Channel, data); err != nil {
		return errors.Wrap(err, "Unable to publish message to Redis")
	}
	if err := conn.Flush(); err != nil {
		return errors.Wrap(err, "Unable to flush published message to Redis")
	}
	return nil
}

// publish to Redis via channel.
func (rw *RedisWriter) publish(data []byte) {
	rw.events <- data
}
