package wshub

// Lobby represents a collection of Rooms.
type Lobby struct {
	Rooms map[string]*Room
}

// NewLobby creates and returns a new Lobby instance.
func NewLobby() *Lobby {
	return &Lobby{
		Rooms: make(map[string]*Room),
	}
}
