package ws

type Message struct {
	Scope string `json:"scope"`
	Data  any    `json:"data"`
}

type ScopeType string

const (
	ScopeTypeWhiteboard       ScopeType = "whiteboard"
	ScopeTypeCursor           ScopeType = "cursor"
	ScopeTypeLobby            ScopeType = "lobby"
	ScopeTypeWhitboardHistory ScopeType = "whiteboard-history"
)
