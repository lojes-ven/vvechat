package ws

import (
	"sync"
)

type Hub struct {
	// 用户ID -> 连接列表 (考虑到多端登录，可以用切片，这里先简单点用单个)
	// 为了简单，我们 assumption: 每个用户只有一个连接，或者发给所有连接
	clients map[uint64]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	rwMutex sync.RWMutex
}

var hubInstance *Hub
var once sync.Once

func GetHub() *Hub {
	once.Do(func() {
		hubInstance = &Hub{
			clients:    make(map[uint64]*Client),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
		go hubInstance.run()
	})
	return hubInstance
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.rwMutex.Lock()
			// 如果已有连接，这里策略是顶掉旧的？或者只保留一个。
			// 简单起见，覆盖。
			if old, ok := h.clients[client.UserID]; ok {
				close(old.Send)
				delete(h.clients, client.UserID)
			}
			h.clients[client.UserID] = client
			h.rwMutex.Unlock()
		case client := <-h.unregister:
			h.rwMutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.rwMutex.Unlock()
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID uint64, message []byte) {
	h.rwMutex.RLock()
	client, ok := h.clients[userID]
	h.rwMutex.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			h.rwMutex.Lock()
			delete(h.clients, userID)
			h.rwMutex.Unlock()
		}
	}
}
