package data

import (
	"strings"
)

const (
	PublicRoom  = iota
	PrivateRoom = iota
)

type ChatRoom struct {
	Title    string `json:"title"`
	User     string `json:"name"`
	Type     uint   `json:"classification"` // 0 = public, 1 = private
	Password string `json:"password"`       // optional
	//	ID       int    `json:"id"`
}

type ChatServer struct {
	//	Rooms        map[int]*ChatRoom
	Rooms map[string]*ChatRoom
	Index *int
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
	//	Rooms:        make(map[int]*ChatRoom),
	Rooms: make(map[string]*ChatRoom),
	Index: &index,
}

func (cs ChatServer) push(cr *ChatRoom) {
	cs.Rooms[strings.ToLower(cr.Title)] = cr
	*cs.Index++
}

func (cs ChatServer) pop(title string) {
	delete(cs.Rooms, strings.ToLower(title))
	*cs.Index--
}

func (cs ChatServer) Chats() (rooms []ChatRoom, err error) {
	rooms = make([]ChatRoom, 0)
	for _, v := range CS.Rooms {
		if v.Type == PublicRoom {
			rooms = append(rooms, *v)
		}
	}
	return
}

// Retrieve returns a single chat room based on title
func Retrieve(title string) (cr ChatRoom, err error) {
	cr = *CS.Rooms[strings.ToLower(title)]
	//err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Create a new chat room
func (cr *ChatRoom) Create() (err error) {
	CS.push(cr)
	//cr.ID = CS.Index // TODO: remove
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
	CS.pop(strings.ToLower(cr.Title))
	//_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}
