package matrix

type RoomEvent struct {
	Content Content `json:"content"`
	EventID string  `json:"event_id"`
	RoomID  string  `json:"room_id"`
	Type    string  `json:"type"`
}
