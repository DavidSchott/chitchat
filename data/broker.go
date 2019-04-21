package data

import (
	"log"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range Clients` started.
const patience time.Duration = time.Second * 1

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	NewClients chan chan []byte

	// Closed client connections
	ClosingClients chan chan []byte

	// Client connections registry
	Clients map[chan []byte]bool
}

func NewBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		NewClients:     make(chan chan []byte),
		ClosingClients: make(chan chan []byte),
		Clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.NewClients:

			// A new client has connected.
			// Register their message channel
			broker.Clients[s] = true
			log.Printf("Client added. %d registered Clients", len(broker.Clients))
			//broker.Notifier <- []byte(fmt.Sprintf("Client added. %d registered Clients", len(broker.Clients)))
		case s := <-broker.ClosingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.Clients, s)
			log.Printf("Removed client. %d registered Clients", len(broker.Clients))
			//broker.Notifier <- []byte(fmt.Sprintf("Removed client. %d registered Clients", len(broker.Clients)))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected Clients
			for clientMessageChan, _ := range broker.Clients {
				select {
				case clientMessageChan <- event:
				case <-time.After(patience):
					log.Print("Skipping client.")
				}
			}
		}
	}

}

/*
func main() {

	broker := NewBroker()

	go func() {
		for {
			time.Sleep(time.Second * 2)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))

}
*/
