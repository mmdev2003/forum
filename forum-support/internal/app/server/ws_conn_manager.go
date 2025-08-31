package server

import (
	"context"
	"forum-support/internal/model"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"time"
)

func NewWsConnManager(
	podID string,
	redisClient model.IRedis,
	messageBroker model.IMessageBroker,
) *WsConnManager {
	return &WsConnManager{
		connections:   make(map[int]*websocket.Conn),
		mu:            &sync.Mutex{},
		redisClient:   redisClient,
		messageBroker: messageBroker,
		podID:         podID,
	}
}

type WsConnManager struct {
	connections   map[int]*websocket.Conn
	mu            *sync.Mutex
	redisClient   model.IRedis
	messageBroker model.IMessageBroker
	podID         string
}

func (m *WsConnManager) AddConnection(ctx context.Context, accountID int, conn *websocket.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.redisClient.Set(ctx, strconv.Itoa(accountID), m.podID, time.Hour*10)
	if err != nil {
		return err
	}
	m.connections[accountID] = conn
	return nil
}

func (m *WsConnManager) RemoveConnection(accountID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, accountID)
}

func (m *WsConnManager) SendMsg(ctx context.Context, toAccountID int, msg []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	conn, ok := m.connections[toAccountID]
	if !ok {
		podID, err := m.redisClient.Get(ctx, strconv.Itoa(toAccountID))
		if err != nil {
			return err
		}

		if podID != m.podID {
			err = m.messageBroker.Publish(ctx, podID, msg)
			return nil
		}
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}
