package matrix

type Content struct {
	Body    string `json:"body"`
	MsgType string `json:"msgtype"`
}

type MessagesResp struct {
	Chunk []RoomEvent `json:"chunk"`
}
