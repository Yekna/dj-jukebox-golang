package websocket

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Manager struct {
	mu        sync.Mutex
	rooms     map[string][]*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string][]*websocket.Conn),
	}
}

func (m *Manager) Connect(room string, c *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rooms[room] = append(m.rooms[room], c)
}

func (m *Manager) Disconnect(room string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conns := m.rooms[room]
	newList := []*websocket.Conn{}

	for _, c := range conns {
		if c != conn {
			newList = append(newList, c)
		}
	}

	if len(newList) == 0 {
		delete(m.rooms, room)
	} else {
		m.rooms[room] = newList
	}
}

func (m *Manager) Broadcast(room string, msg interface{}) {
	m.mu.Lock()
	conns := m.rooms[room]
	m.mu.Unlock()

	for _, c := range conns {
		if err := c.WriteJSON(msg); err != nil {
			// ignore
		}
	}
}

