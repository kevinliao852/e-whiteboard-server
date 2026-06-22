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

func (m *fakeCanvasModel) GetByWhiteboardID(whiteboardID int) ([]core.CanvasData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]core.CanvasData, 0)
	for _, entry := range m.data {
		if entry.WhiteboardId == whiteboardID {
			result = append(result, *entry)
		}
	}

	return result, nil
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

func TestDrawingSVC_ListCanvasData(t *testing.T) {
	model := &fakeCanvasModel{
		data: []*core.CanvasData{
			{ID: 1, WhiteboardId: 7, StartX: 10, StartY: 20, EndX: 30, EndY: 40},
			{ID: 2, WhiteboardId: 8, StartX: 11, StartY: 21, EndX: 31, EndY: 41},
			{ID: 3, WhiteboardId: 7, StartX: 12, StartY: 22, EndX: 32, EndY: 42},
		},
	}
	svc := NewDrawingSVC(model)

	got, err := svc.ListCanvasData(7)
	if err != nil {
		t.Fatalf("ListCanvasData returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 points, got %d", len(got))
	}
	if got[0].WhiteboardId != 7 || got[1].WhiteboardId != 7 {
		t.Fatalf("expected only whiteboard 7 data, got %#v", got)
	}
}
