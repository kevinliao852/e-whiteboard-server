package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type drawingRoomServiceStub struct{}

func (s *drawingRoomServiceStub) CreateRoom(roomID string) (*core.Room, error) { return nil, nil }
func (s *drawingRoomServiceStub) JoinRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *drawingRoomServiceStub) LeaveRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *drawingRoomServiceStub) BroadcastToRoom(roomID string, message string) error { return nil }
func (s *drawingRoomServiceStub) ListRooms() []core.Room                              { return nil }

type drawingServiceStub struct {
	points []core.CanvasData
}

func (s *drawingServiceStub) Enqueue(roomID string, message []byte) error { return nil }

func (s *drawingServiceStub) ListCanvasData(whiteboardID int) ([]core.CanvasData, error) {
	return s.points, nil
}

func TestDrawingController_GetPoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := DrawingController{
		RoomService: &drawingRoomServiceStub{},
		DrawingService: &drawingServiceStub{
			points: []core.CanvasData{
				{ID: 1, WhiteboardId: 7, StartX: 10, StartY: 20, EndX: 30, EndY: 40},
				{ID: 2, WhiteboardId: 7, StartX: 12, StartY: 22, EndX: 32, EndY: 42},
			},
		},
	}

	router := gin.Default()
	router.GET("/whiteboards/:id/points", ctrl.GetPoints)

	t.Run("returns points", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/whiteboards/7/points", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response []PointResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Fatalf("expected 2 points, got %d", len(response))
		}
		if response[0].Start != [2]int{10, 20} {
			t.Fatalf("unexpected first start point: %#v", response[0].Start)
		}
		if response[1].End != [2]int{32, 42} {
			t.Fatalf("unexpected second end point: %#v", response[1].End)
		}
	})

	t.Run("rejects invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/whiteboards/abc/points", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})
}

func TestDrawingController_DrawRejectsInvalidWhiteboardID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := DrawingController{
		RoomService:    &drawingRoomServiceStub{},
		DrawingService: &drawingServiceStub{},
	}

	router := gin.Default()
	router.GET("/drawing/:id", ctrl.Draw())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/drawing/not-a-number", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
