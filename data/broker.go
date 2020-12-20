package data

import (
	"log"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range Clients` started.
const patience time.Duration = time.Second * 1

// Broker maintains the client connections and handles events using a listener goroutine
type Broker struct {
	// Registered Clients.
	Clients map[*Client]bool

	// Inbound messages from the Clients.
	Notification chan []byte

	// Register requests from the Clients.
	OpenClient chan *Client

	// Unregister requests from Clients.
	CloseClient chan *Client

	RoomID int
}

func newBroker(ID int) *Broker {
	return &Broker{
		Notification: make(chan []byte),
		OpenClient:   make(chan *Client),
		CloseClient:  make(chan *Client),
		Clients:      make(map[*Client]bool),
		RoomID:       ID,
	}
}

func (br *Broker) listen() {
	for {
		select {
		case c := <-br.OpenClient:
			// A new client has connected.
			// Register their message channel
			br.Clients[c] = true
			log.Printf("Client added. %d registered Clients", len(br.Clients))
		case c := <-br.CloseClient:
			// A client has dettached and we want to
			// stop sending them messages.
			if _, ok := br.Clients[c]; ok {
				delete(br.Clients, c)
				close(c.Send)
				//cr, _ := CS.Retrieve(strconv.Itoa(br.RoomID))
				//c.unsubscribe(&ChatEvent{User: c.Username})
				/*if err := cr.RemoveClient(c.Username); err != nil {
					log.Printf("Error removing client: %s from room %s. Error: %s", c.Username, cr.Title, err.Error())
				}*/
				log.Printf("Removed client. %d registered Clients", len(br.Clients))
			}
		case evt := <-br.Notification:
			// We got a new event from the outside
			// Send event to all connected Clients
			for client := range br.Clients {
				select {
				case client.Send <- evt:
				case <-time.After(patience):
					log.Print("Skipping client: " + client.Username)
				default:
					log.Print("Deleting client: " + client.Username)
					close(client.Send)
					delete(br.Clients, client)
				}
			}
		}
	}
}
