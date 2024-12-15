// main.go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/redis/go-redis/v9"
    "github.com/gorilla/websocket"
    "golang.org/x/net/context"
)

type Server struct {
    redisClient *redis.Client
    upgrader    websocket.Upgrader
    // Map to store active connections
    clients     map[*websocket.Conn]string // websocket -> channelID
    clientsLock sync.RWMutex
}

type Message struct {
    Type      string `json:"type"`
    ChannelID string `json:"channelId"`
    UserID    string `json:"userId"`
    Action    string `json:"action"`
}

func NewServer() *Server {
    // Initialize Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })

    return &Server{
        redisClient: rdb,
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true // In production, implement proper origin checking
            },
        },
        clients: make(map[*websocket.Conn]string),
    }
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()

    // Handle WebSocket connection
    for {
        var msg Message
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        switch msg.Type {
        case "join_channel":
            s.handleJoinChannel(conn, msg)
        case "user_action":
            s.handleUserAction(msg)
        }
    }
    s.clientsLock.Lock()
    delete(s.clients, conn)
    s.clientsLock.Unlock()
}

func (s *Server) handleJoinChannel(conn *websocket.Conn, msg Message) {
    s.clientsLock.Lock()
    s.clients[conn] = msg.ChannelID
    s.clientsLock.Unlock()
    ctx := context.Background()
    s.redisClient.SAdd(ctx, "channel:"+msg.ChannelID+":users", msg.UserID)

    // Subscribe to Redis channel
    pubsub := s.redisClient.Subscribe(ctx, "channel:"+msg.ChannelID)
    go func() {
        for {
            msg, err := pubsub.ReceiveMessage(ctx)
            if err != nil {
                log.Printf("Redis subscription error: %v", err)
                return
            }
            // Broadcast message to all clients in channel
            s.broadcastToChannel(msg.Channel, msg.Payload)
        }
    }()
}

func (s *Server) handleUserAction(msg Message) {
    // Publish action to Redis channel
    ctx := context.Background()
    payload, _ := json.Marshal(msg)
    s.redisClient.Publish(ctx, "channel:"+msg.ChannelID, payload)
}

func (s *Server) broadcastToChannel(channelID string, payload string) {
    s.clientsLock.RLock()
    defer s.clientsLock.RUnlock()

    for conn, clientChannel := range s.clients {
        if clientChannel == channelID {
            err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
            if err != nil {
                log.Printf("Error broadcasting message: %v", err)
            }
        }
    }
}

func main() {
    server := NewServer()

    // Define routes
    http.HandleFunc("/ws", server.handleWebSocket)

    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}