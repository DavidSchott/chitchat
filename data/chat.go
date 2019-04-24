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

type ChatRoom struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Type        string             `json:"classification"`
	Password    string             `json:"password"` // optional
	CreatedAt   time.Time          `json:"time"`
	ID          int                `json:"id"`
	Broker      *Broker            `json:"-"`
	Clients     map[string]*Client `json:"-"`
}

type ChatServer struct {
	RoomsID map[int]*ChatRoom
	Rooms   map[string]*ChatRoom
	Index   *int
}

type Success struct {
	Sucess bool `json:"sucess"`
}

type Failure struct {
	Sucess bool   `json:"sucess"`
	Error  string `json:"error"`
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
			code: 101,
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
		for k, _ := range CS.Rooms {
			if k == titleorID {
				return true
			}
		}
	}
	return false
}

// Create a new chat room
func (cs ChatServer) Add(cr *ChatRoom) (err error) {
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
