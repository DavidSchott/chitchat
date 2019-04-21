package data

import (
	"strings"
	"time"
)

const (
	PublicRoom  = "public"
	PrivateRoom = "private"
	HiddenRoom  = "hidden"
)

type ChatRoom struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	User         string    `json:"name"`
	Type         string    `json:"classification"` // 0 = public, 1 = private
	Password     string    `json:"password"`       // optional
	CreatedAt    time.Time `json:"time"`
	Participants []string  `json:"participants"`
	ID           int       `json:"id"`
	Broker       *Broker
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
		Title:        "Public Chat",
		Description:  "This is the default chat, available to everyone!",
		User:         "Server",
		Type:         "public",
		Password:     "",
		CreatedAt:    time.Now(),
		Participants: []string{"Server"},
		ID:           0,
		Broker:       NewBroker(),
	})
}

func (cs ChatServer) push(cr *ChatRoom) {
	// Update indices, create new session
	*cs.Index++
	cr.ID = *cs.Index
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

// Retrieve returns a single chat room based on title
func Retrieve(title string) (cr ChatRoom, err error) {
	cr = *CS.Rooms[strings.ToLower(title)]
	//err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Retrieve returns a single chat room based on ID
func RetrieveID(ID int) (cr ChatRoom, err error) {
	cr = *CS.RoomsID[ID]
	//err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Create a new chat room
func (cr *ChatRoom) Create() (err error) {
	cr.CreatedAt = time.Now()
	cr.Participants = []string{cr.User}
	cr.Broker = NewBroker()
	CS.push(cr)
	return
}

// Update a chat room
func (cr *ChatRoom) Update() (err error) {
	CS.Rooms[strings.ToLower(cr.Title)] = cr
	//_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

// Delete a chat room
func (cr *ChatRoom) Delete() (err error) {
	CS.pop(strings.ToLower(cr.Title), cr.ID)
	//_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}
