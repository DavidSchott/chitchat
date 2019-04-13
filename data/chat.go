package data

type ChatRoom struct {
	Title    string
	User     string
	Type     uint   // 0 = public, 1 = private
	Password string // optional
	ID       int
}

type ChatServer struct {
	Rooms        map[int]*ChatRoom
	RoomsByTitle map[string]*ChatRoom
	Index        int
}

func (cs ChatServer) Push(cr *ChatRoom) {
	cs.Rooms[cs.Index] = cr
	cs.RoomsByTitle[cr.Title] = cr
}

func (cs ChatServer) Pop(ID int, title string) {
	cs.Rooms[ID] = nil
	cs.RoomsByTitle[title] = nil
}
