package downloads

import (
	"sync"

	"ytdl-api/internal/models"
)

const subscriberBuffer = 512

type ProgressEvent struct {
	PlaylistID   int64                 `json:"playlist_id"`
	EntryID      int64                 `json:"entry_id"`
	Status       models.DownloadStatus `json:"status"`
	FilePath     string                `json:"file_path,omitempty"`
	ErrorMessage string                `json:"error_message,omitempty"`
	ErrorCode    models.EntryErrorCode `json:"error_code,omitempty"`
	Phase        models.DownloadPhase  `json:"phase,omitempty"`
	Percent      *float64              `json:"percent,omitempty"`
	ETA          *int64                `json:"eta,omitempty"`
}

func (ev ProgressEvent) droppable() bool {
	return ev.Phase != ""
}

type Hub struct {
	mu          sync.Mutex
	subscribers map[int64]map[chan ProgressEvent]struct{}
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[int64]map[chan ProgressEvent]struct{}),
	}
}

func (h *Hub) Subscribe(playlistID int64) chan ProgressEvent {
	ch := make(chan ProgressEvent, subscriberBuffer)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscribers[playlistID] == nil {
		h.subscribers[playlistID] = make(map[chan ProgressEvent]struct{})
	}
	h.subscribers[playlistID][ch] = struct{}{}
	return ch
}

func (h *Hub) Unsubscribe(playlistID int64, ch chan ProgressEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.remove(playlistID, ch)
}

func (h *Hub) Broadcast(playlistID int64, ev ProgressEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subscribers[playlistID] {
		select {
		case ch <- ev:
		default:
			if !ev.droppable() {
				h.remove(playlistID, ch)
			}
		}
	}
}

func (h *Hub) remove(playlistID int64, ch chan ProgressEvent) {
	subs, ok := h.subscribers[playlistID]
	if !ok {
		return
	}
	if _, ok := subs[ch]; !ok {
		return
	}
	delete(subs, ch)
	if len(subs) == 0 {
		delete(h.subscribers, playlistID)
	}
	close(ch)
}
