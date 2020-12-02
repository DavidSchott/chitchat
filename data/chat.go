package data

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const (
	PublicRoom  = "public"
	PrivateRoom = "private"
	HiddenRoom  = "hidden"
)

type Outcome struct {
	Status bool      `json:"status"`
	Error  *APIError `json:"error,omitempty"`
}

// ChatRoom is a struct representing a chat room
// TODO:  Add Administrator
type ChatRoom struct {
	Title       string             `json:"title"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"visibility"`
	Password    string             `json:"password,omitempty"` // TODO: Make this json:- once salted
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	ID          int                `json:"id"`
	Broker      *Broker            `json:"-"`
	Clients     map[string]*Client `json:"-"`
}

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
		Clients []Client `json:"participants"`
	}{
		ChatRoom: &cr,
		Clients:  clientsSlice,
	})
	return jsonEncoding, err
}

func (cr ChatRoom) AddClient(c *Client) (err error) {
	if cr.clientExists(c.Username) {
		return &APIError{
			Code:  202,
			Field: c.Username,
		}
	}
	cr.Clients[strings.ToLower(c.Username)] = c
	return
}

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
	if (len(cr.Password) < 8 || len(cr.Password) > 20) && visibility != PublicRoom {
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
	return cr.Password == val // TODO: Salted passwords
}

func (cr ChatRoom) clientExists(name string) bool {
	name = strings.ToLower(name)
	for k, _ := range cr.Clients {
		if k == name {
			return true
		}
	}
	return false
}

type ChatServer struct {
	RoomsID map[int]*ChatRoom
	Rooms   map[string]*ChatRoom // TODO: Remove this duplication once data layer moves to DB
	Index   *int
}

var index int
var CS ChatServer = ChatServer{
	RoomsID: make(map[int]*ChatRoom),
	Rooms:   make(map[string]*ChatRoom),
	Index:   &index,
}

func (cs ChatServer) Init() {
	CS.push(&ChatRoom{
		Title:       "Public Chat",
		Description: "This is the default chat, available to everyone!",
		Type:        "public",
		Password:    "",
		CreatedAt:   time.Now(),
		ID:          0,
		Broker:      NewBroker(),
		Clients:     make(map[string]*Client),
	})
}

func (cs ChatServer) push(cr *ChatRoom) {
	// Update indices, create new session
	*cs.Index++
	cr.ID = *cs.Index
	cr.Clients = make(map[string]*Client)
	cr.Type = strings.ToLower(cr.Type)
	// Push to chat server
	cs.Rooms[strings.ToLower(cr.Title)] = cr
	cs.RoomsID[cr.ID] = cr
}

func (cs ChatServer) pop(title string, ID int) {
	delete(cs.Rooms, strings.ToLower(title))
	delete(cs.RoomsID, ID)
	*cs.Index--
}

func (cs ChatServer) Chats() (rooms []ChatRoom, err error) {
	rooms = make([]ChatRoom, 0)
	for _, v := range CS.Rooms {
		if v.Type != HiddenRoom {
			rooms = append(rooms, *v)
		}
	}
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

// Retrieve returns a single chat room based on title or ID
func (cs ChatServer) Retrieve(title string) (cr *ChatRoom, err error) {
	if !cs.roomExists(title) {
		return cr, &APIError{
			Code:  101,
			Field: title,
		}
	}
	if id := isInt(title); id != -1 {
		cr = cs.RoomsID[id]
	} else {
		cr = cs.Rooms[strings.ToLower(title)]
	}
	//err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return cr, nil
}

// RetrieveID returns a single chat room based on ID. NOTE: This has no error handling unlike cs.Retrieve()
func (cs ChatServer) RetrieveID(ID int) (cr *ChatRoom, err error) {
	cr = cs.RoomsID[ID]
	//err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

func (cs ChatServer) roomExists(titleorID string) bool {
	if id, err := strconv.Atoi(titleorID); err == nil {
		for k, _ := range CS.RoomsID {
			if k == id {
				return true
			}
		}
	} else {
		titleorID = strings.ToLower(titleorID)
		for k, _ := range CS.Rooms {
			if strings.ToLower(k) == titleorID {
				return true
			}
		}
	}
	return false
}

// Add will create a new chat room and add it to the server
func (cs ChatServer) Add(cr *ChatRoom) (err error) {
	// validate chat room request
	if apierr, valid := cr.IsValid(); !valid {
		return apierr
	}
	if cs.roomExists(cr.Title) { // TODO: What if the room is hidden? Return unspecified error or inform user?
		return &APIError{
			Code:  102,
			Field: cr.Title,
		}
	}
	cr.CreatedAt = time.Now()
	cr.Broker = NewBroker()
	cr.Type = strings.ToLower(cr.Type)
	cs.push(cr)
	return
}

// Update a chat room. NOTE: Authorization should have been done before calling this
// TODO: Get input from requested ID. Edit both RoomsID and Rooms.
func (cs ChatServer) Update(titleOrID string, modifiedChatRoom *ChatRoom) (err error) {
	currentChatRoom, err := cs.Retrieve(titleOrID)
	if err != nil {
		return
	}
	// Update password for validation
	modifiedChatRoom.Password = currentChatRoom.Password
	if apierr, valid := modifiedChatRoom.IsValid(); !valid {
		return apierr
	}

	/* 	This is commented for now since modifying a password or visibility could be a legitimate use-case. Can authorize using cookie
	Check password matches.
	if cr.Type != PublicRoom && !cs.RoomsID[cr.ID].MatchesPassword(cr.Password) {
		return &APIError{
			Code:  104,
			Field: "password",
		}
	}
	// Ensure room type is not trying to be changed.
	if cr.Type != cs.RoomsID[cr.ID].Type {
		return &APIError{
			Code:  104,
			Field: "visibility",
		}
	}
	*/
	// Update chat room
	// TODO: Ensure ID is not modified, update time
	modifiedChatRoom.ID = currentChatRoom.ID
	modifiedChatRoom.UpdatedAt = time.Now()
	*currentChatRoom = *modifiedChatRoom
	//_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

// Delete a chat room
func (cs ChatServer) Delete(cr *ChatRoom) (err error) {
	cs.pop(strings.ToLower(cr.Title), cr.ID)
	//_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}
