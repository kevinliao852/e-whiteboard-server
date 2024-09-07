package wshub

type Lobby struct {
	Rooms map[string]*Room
}

func NewLobby() *Lobby {
	return &Lobby{
		Rooms: make(map[string]*Room),
	}
}
