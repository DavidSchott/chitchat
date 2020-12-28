package data

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// PublicRoom is a room open for anyone to join without authentication
	PublicRoom = "public"
	// PrivateRoom is password protected and requires an authentication token in order to process requests
	PrivateRoom = "private"
	// HiddenRoom is a private room that is not listed on public-facing APIs. TODO: Hide this from GET /chats/<id> as well?
	HiddenRoom = "hidden"
)

// ChatRoom is a struct representing a chat room
// TODO:  Add Administrator
type ChatRoom struct {
	Title       string             `json:"title"`
	ID          string             `json:"id"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"visibility"`
	Password    string             `json:"password,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Broker      *Broker            `json:"-"`
	Clients     map[string]*Client `json:"-"`
}

// ToJSON marshals a ChatRoom object in a JSON encoding that can be returned to users
func (cr ChatRoom) ToJSON() (jsonEncoding []byte, err error) {
	// Populate client slice. TODO: Can this be simplified?
	clientsSlice := make([]Client, len(cr.Clients))
	var i int = 0
	for _, v := range cr.Clients {
		//clientsSlice = append(clientsSlice, *v)
		clientsSlice[i] = *v
		i++
	}
	// Create new JSON struct with clients
	jsonEncoding, err = json.Marshal(struct {
		*ChatRoom
		Clients []Client `json:"users"`
	}{
		ChatRoom: &cr,
		Clients:  clientsSlice,
	})
	return jsonEncoding, err
}

//AddClient will add a user to a ChatRoom
func (cr ChatRoom) AddClient(c *Client) (err error) {
	if cr.clientExists(c.Username) {
		return &APIError{
			Code:  202,
			Field: c.Username,
		}
	}
	c.LastActivity = time.Now()
	c.UserID = generateUUID()
	cr.Clients[strings.ToLower(c.Username)] = c
	return
}

// RemoveClient will remove a user from a ChatRoom
func (cr ChatRoom) RemoveClient(user string) (err error) {
	if !cr.clientExists(user) {
		return &APIError{
			Code:  201,
			Field: user,
		}
	}
	delete(cr.Clients, strings.ToLower(user))
	return
}

// Authorize authorizes a given ChatEvent for the Room
func (cr ChatRoom) Authorize(c *ChatEvent) bool {
	return cr.MatchesPassword(c.Password)
}

// IsValid validates a chat room fields are still valid
func (cr ChatRoom) IsValid() (err *APIError, validity bool) {
	// Title should be at least 2 characters
	if len(cr.Title) < 2 || len(cr.Title) > 70 {
		return &APIError{
			Code:  105,
			Field: "title",
		}, false
	}
	// Description shall not be too long
	if len(cr.Description) > 70 {
		return &APIError{
			Code:  105,
			Field: "description",
		}, false
	}
	visibility := strings.ToLower(cr.Type)
	// Visibility must be set
	if visibility != PublicRoom && visibility != PrivateRoom && visibility != HiddenRoom {
		return &APIError{
			Code:  105,
			Field: "visibility",
		}, false
	}
	// Non-public rooms require a valid password
	if (len(cr.Password) < 8) && visibility != PublicRoom {
		return &APIError{
			Code:  105,
			Field: "password",
		}, false
	}
	// A public room should not have a password set (to avoid accidents)
	if len(cr.Password) != 0 && visibility == PublicRoom {
		return &APIError{
			Code:  105,
			Field: "visibility",
		}, false
	}
	return nil, true
}

// MatchesPassword takes in a value and compares it with the room's password
func (cr ChatRoom) MatchesPassword(val string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cr.Password), []byte(val))
	return err == nil
}

func (cr ChatRoom) clientExists(name string) bool {
	name = strings.ToLower(name)
	for k := range cr.Clients {
		if k == name {
			return true
		}
	}
	return false
}

// create a random UUID (adheres to RFC 4122)
// adapted from http://github.com/nu7hatch/gouuid
func generateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// PrettyTime prints the creation date in a pretty format
func (cr ChatRoom) PrettyTime() string {
	layout := "Mon Jan _2 15:04"
	return cr.CreatedAt.Format(layout)
}

// Participants prints the # of active clients
func (cr ChatRoom) Participants() int {
	return len(cr.Clients)
}
