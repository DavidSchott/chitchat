package data

import (
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
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ChatRoom struct {
	Title       string             `json:"title"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"classification"`
	Password    string             `json:"password,omitempty"` // optional
	CreatedAt   time.Time          `json:"time"`
	ID          int                `json:"id"`
	Broker      *Broker            `json:"-"`
	Clients     map[string]*Client `json:"-"`
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

// Authenticate authenticates a given ChatEvent for the Room
func (cr ChatRoom) Authenticate(c *ChatEvent) bool {
	return cr.MatchesPassword(c.Password)
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
	Rooms   map[string]*ChatRoom
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
	CS.push(&ChatRoom{
		Title:       "Private Chat",
		Description: "This is a test private chat, used for testing!",
		Type:        "private",
		Password:    "secret123",
		CreatedAt:   time.Now(),
		ID:          99,
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

// Retrieve returns a single chat room based on title
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
	return
}

// Retrieve returns a single chat room based on ID
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

// Create a new chat room
func (cs ChatServer) Add(cr *ChatRoom) (err error) {
	//fmt.Println(cs.roomExists(cr.Title))
	if cs.roomExists(cr.Title) { // TODO: What if the room is hidden?
		return &APIError{
			Code:  102,
			Field: cr.Title,
		}
	}
	cr.CreatedAt = time.Now()
	cr.Broker = NewBroker()
	cs.push(cr)
	return
}

// Update a chat room
func (cs ChatServer) Update(cr *ChatRoom) (err error) {
	cs.Rooms[strings.ToLower(cr.Title)] = cr
	//_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

// Delete a chat room
func (cs ChatServer) Delete(cr *ChatRoom) (err error) {
	cs.pop(strings.ToLower(cr.Title), cr.ID)
	//_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}
