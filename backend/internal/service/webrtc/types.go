package webrtc

import (
	"context"
	"github.com/gorilla/websocket"
	"time"
)

type SignalingMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Session struct {
	sessionId       string
	prompt          string
	closed          bool
	ctx             context.Context
	cancel          context.CancelFunc
	frontendConn    *websocket.Conn
	pbxConn         *websocket.Conn
	lastInterruptAt time.Time
	ttsPlaying      bool
}
