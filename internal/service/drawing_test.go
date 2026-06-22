package service

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type fakeCanvasModel struct {
	mu   sync.Mutex
	data []*core.CanvasData
	ch   chan *core.CanvasData
}

func (m *fakeCanvasModel) Create(data *core.CanvasData) error {
	m.mu.Lock()
	m.data = append(m.data, data)
	m.mu.Unlock()

	select {
	case m.ch <- data:
	default:
	}
	return nil
}

func TestDrawingSVC_EnqueueWhiteboardMessage(t *testing.T) {
	model := &fakeCanvasModel{ch: make(chan *core.CanvasData, 1)}
	svc := NewDrawingSVC(model)

	message := map[string]any{
		"scope": "whiteboard",
		"data": map[string]any{
			"start": []int{10, 20},
			"end":   []int{30, 40},
		},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	if err := svc.Enqueue("123", payload); err != nil {
		t.Fatalf("Enqueue returned error: %v", err)
	}

	select {
	case got := <-model.ch:
		if got.WhiteboardId != 123 || got.StartX != 10 || got.StartY != 20 || got.EndX != 30 || got.EndY != 40 {
			t.Fatalf("unexpected canvas data: %#v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for canvas data to be persisted")
	}
}

func TestDrawingSVC_IgnoreNonWhiteboardMessage(t *testing.T) {
	model := &fakeCanvasModel{ch: make(chan *core.CanvasData, 1)}
	svc := NewDrawingSVC(model)

	message := map[string]any{
		"scope": "lobby",
		"data": map[string]any{
			"room_id": "123",
		},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	if err := svc.Enqueue("123", payload); err != nil {
		t.Fatalf("Enqueue returned error: %v", err)
	}

	select {
	case got := <-model.ch:
		t.Fatalf("expected no persistence for lobby message, got %#v", got)
	case <-time.After(200 * time.Millisecond):
	}
}
