package wshub

type Message struct {
	Scope string `json:"scope"`
	Data  any    `json:"data"`
}

type ScopeType string

const (
	ScopeTypeWhiteboard       ScopeType = "whiteboard"
	ScopeTypeLobby            ScopeType = "lobby"
	ScopeTypeWhitboardHistory ScopeType = "whiteboard-history"
)
