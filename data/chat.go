package data

import "strings"

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
	Index int
}

type Success struct {
	Sucess bool `json:"sucess"`
}

type Failure struct {
	Sucess bool   `json:"sucess"`
	Error  string `json:"error"`
}

var CS ChatServer = ChatServer{
	//	Rooms:        make(map[int]*ChatRoom),
	Rooms: make(map[string]*ChatRoom),
	Index: 0,
}

func (cs ChatServer) push(cr *ChatRoom) {
	cs.Rooms[strings.ToLower(cr.Title)] = cr
}

func (cs ChatServer) pop(title string) {
	cs.Rooms[strings.ToLower(title)] = nil
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
	CS.Index++
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
