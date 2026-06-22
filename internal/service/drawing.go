package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type DrawingSVC struct {
	CreateFn            func(data *core.CanvasData) error
	GetByWhiteboardIDFn func(whiteboardID int) ([]core.CanvasData, error)
	queue               chan *core.CanvasData
	once                sync.Once
}

func NewDrawingSVC(createFn func(data *core.CanvasData) error, getByWhiteboardIDFn func(whiteboardID int) ([]core.CanvasData, error)) *DrawingSVC {
	svc := &DrawingSVC{
		CreateFn:            createFn,
		GetByWhiteboardIDFn: getByWhiteboardIDFn,
		queue:               make(chan *core.CanvasData, 100),
	}
	svc.startWorker()
	return svc
}

func (s *DrawingSVC) startWorker() {
	s.once.Do(func() {
		go func() {
			for data := range s.queue {
				if data == nil || s.CreateFn == nil {
					continue
				}
				if err := s.CreateFn(data); err != nil {
					log.Printf("failed to persist drawing data: %v", err)
				}
			}
		}()
	})
}

type incomingMessage struct {
	Scope string          `json:"scope"`
	Data  json.RawMessage `json:"data"`
}

type drawingPayload struct {
	Start [2]int `json:"start"`
	End   [2]int `json:"end"`
}

func (s *DrawingSVC) Enqueue(roomID string, message []byte) error {
	if s == nil {
		return nil
	}

	var parsed incomingMessage
	if err := json.Unmarshal(message, &parsed); err != nil {
		return err
	}

	if parsed.Scope != "whiteboard" {
		return nil
	}

	data, err := s.parseCanvasData(roomID, parsed.Data)
	if err != nil {
		return err
	}

	select {
	case s.queue <- data:
	default:
		log.Printf("dropping drawing update for room %s: queue full", roomID)
	}

	return nil
}

func (s *DrawingSVC) ListCanvasData(whiteboardID int) ([]core.CanvasData, error) {
	if s == nil || s.GetByWhiteboardIDFn == nil {
		return nil, nil
	}

	return s.GetByWhiteboardIDFn(whiteboardID)
}

func (s *DrawingSVC) parseCanvasData(roomID string, raw json.RawMessage) (*core.CanvasData, error) {
	whiteboardID, err := strconv.Atoi(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid room id %q: %w", roomID, err)
	}

	var payload drawingPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("invalid drawing payload: %w", err)
	}

	return &core.CanvasData{
		WhiteboardId: whiteboardID,
		StartX:       payload.Start[0],
		StartY:       payload.Start[1],
		EndX:         payload.End[0],
		EndY:         payload.End[1],
	}, nil
}

var _ core.DrawingService = (*DrawingSVC)(nil)
