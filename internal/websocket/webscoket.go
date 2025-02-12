package websocket

import (
	"log"
	"net/http"
	"sync"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket upgrader, to upgrade http to websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client represents a connected users
type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID string
}

// Room to manage to clients
type Room struct {
	Clients map[string]*Client
	Mutex   sync.Mutex
}

// room manager , manages multiple rooms
type RoomManager struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

var roomManager = &RoomManager{
	Rooms: make(map[string]*Room),
}

func HandleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println("Error upgrading to websocket ")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error upgrading to websocket " + err.Error()})
		return
	}

	clientID := c.Query("userId")
	roomID := c.Query("roomID")

	client := &Client{
		ID:     clientID,
		Conn:   conn,
		Send:   make(chan []byte),
		RoomID: roomID,
	}

	go client.readMessage()
	go client.writeMessage()

	// add client to the room
	addClientToRoom(roomID, client)
}

func (c *Client) readMessage() {
	defer func() {
		fmt.Println("Client removed ", c.ID)
		removeClientFromRoom(c.RoomID, c.ID)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Message read error")
			break
		}
		broadcastToRoom(c.RoomID, message, c.ID)
	}

}

func (c *Client) writeMessage() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message")
			break
		}
	}
}

func addClientToRoom(roomID string, client *Client) {
	roomManager.Mutex.Lock()
	defer roomManager.Mutex.Unlock()

	if _, exists := roomManager.Rooms[roomID]; !exists {
		log.Println("Room with ID " + roomID + " does not exists, creating new Room")
		roomManager.Rooms[roomID] = &Room{
			Clients: make(map[string]*Client),
		}
	}

	room := roomManager.Rooms[roomID]
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	room.Clients[client.ID] = client

	log.Println("Client added to room " + client.ID + " Room: " + roomID)
}

func removeClientFromRoom(roomID string, clientID string) {
	roomManager.Mutex.Lock()
	defer roomManager.Mutex.Unlock()

	room, exists := roomManager.Rooms[roomID]
	if !exists {
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	delete(room.Clients, clientID)

	log.Println("Client " + clientID + " Removed from roomID " + roomID)
}

func broadcastToRoom(roomID string, message []byte, clientID string) {
	roomManager.Mutex.Lock()
	defer roomManager.Mutex.Unlock()
	room, exists := roomManager.Rooms[roomID]
	if !exists {
		fmt.Println("Room does not exists")
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	for _, client := range room.Clients {
		if client.ID == clientID {
			continue
		}
		client.Send <- message
	}
}
