package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

func NewWsConnManager() *WsConnManager {
	return &WsConnManager{
		connections: make(map[int]*websocket.Conn),
		mu:          &sync.Mutex{},
	}
}

type WsConnManager struct {
	connections map[int]*websocket.Conn
	mu          *sync.Mutex
}

func (m *WsConnManager) AddConnection(accountID int, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[accountID] = conn
}

func (m *WsConnManager) RemoveConnection(accountID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, accountID)
}

func (m *WsConnManager) SendMsg(toAccountID int, msg []byte) error {
	m.mu.Lock()
	conn, ok := m.connections[toAccountID]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("not found connection with id %d", toAccountID)
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}
