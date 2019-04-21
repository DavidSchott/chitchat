package data

type Client struct {
	Username string
	Color    string
	Conn     chan chan []byte
}
