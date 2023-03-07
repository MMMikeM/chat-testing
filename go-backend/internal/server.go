package internal

import (
	"context"
	"encoding/json"
	"messanger/internal/openapi"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type connection struct {
	ws             *websocket.Conn
	conversationId string
}

// pongHandler is used to handle PongMessages for the Client
func (c *connection) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	return c.ws.SetReadDeadline(time.Now().Add(pongWait))
}

type conversationRoom struct {
	connections map[*websocket.Conn]bool
	sync.RWMutex
}

func (cr *conversationRoom) addConnection(c *websocket.Conn) {
	cr.Lock()
	defer cr.Unlock()

	cr.connections[c] = true
}

func (cr *conversationRoom) disconnectConnection(c *websocket.Conn) {
	cr.Lock()
	defer cr.Unlock()
	delete(cr.connections, c)
	c.Close()
}

type Server struct {
	Logger              *logrus.Logger
	DB                  *gorm.DB
	ActiveConversations map[string]*conversationRoom
	Connections         chan connection
	Disconnects         chan connection
	NewMessagesChan     chan Message
}

var (
	upgrader        = websocket.Upgrader{}
	pongWait        = 10 * time.Second
	pingInterval    = (pongWait * 9) / 10
	pingPongEnabled = false
)

func NewServer(ctx context.Context, logger *logrus.Logger, DB *gorm.DB) *Server {
	s := Server{
		Logger:              logger,
		DB:                  DB,
		ActiveConversations: map[string]*conversationRoom{},
		Connections:         make(chan connection),
		Disconnects:         make(chan connection),
		NewMessagesChan:     make(chan Message),
	}

	go func() {
		ticker := time.NewTicker(pingInterval)
		for {
			select {
			case msg, ok := <-s.NewMessagesChan:
				if !ok {
					s.Logger.Println("Error reading message from channel")
					// convo := s.ActiveConversations[msg.ConversationId]
					// if err := s.ActiveConversations[msg.ConversationId].WriteMessage(websocket.CloseMessage, nil); err != nil {
					// 	log.Println("connection closed: ", err)
					// }
					return
				}

				jsonMessage, err := json.Marshal(msg)
				if err != nil {
					s.Logger.Println(err)
				}

				for ws := range s.ActiveConversations[msg.ConversationId].connections {
					if err := ws.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
						s.Logger.Println(err)
					}
				}
			case <-ticker.C:
				// 	// Send the Ping
				if pingPongEnabled {
					for _, convoRoom := range s.ActiveConversations {
						for ws := range convoRoom.connections {
							if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
								return // return to break this goroutine triggeing cleanup
							}
						}
					}
				}
			}
		}
	}()

	go func() {
		for conn := range s.Connections {
			if s.ActiveConversations[conn.conversationId] == nil {
				convoRoom := conversationRoom{connections: map[*websocket.Conn]bool{}}
				s.ActiveConversations[conn.conversationId] = &convoRoom
				convoRoom.addConnection(conn.ws)
			} else {
				convoRoom := s.ActiveConversations[conn.conversationId]
				convoRoom.addConnection(conn.ws)
			}
		}
	}()

	go func() {
		for conn := range s.Disconnects {
			convoRoom := s.ActiveConversations[conn.conversationId]
			convoRoom.disconnectConnection(conn.ws)
		}
	}()

	go func() {
		for {
			totalCount := 0
			for _, conv := range s.ActiveConversations {
				totalCount += len(conv.connections)
			}

			s.Logger.Printf("There is currently: %d WS connection(s)", totalCount)
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	return &s
}

func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, "Running")
}

func (s *Server) CreateConversation(ctx echo.Context) error {
	id := uuid.New().String()
	c := Conversation{ID: id, CreatedAt: time.Now()}
	result := s.DB.Create(&c)
	if result.Error != nil {
		return result.Error
	}

	return ctx.JSON(http.StatusCreated, c)
}

func (s *Server) GetConversation(ctx echo.Context, conversationId string) error {
	c := Conversation{}
	result := s.DB.
		First(&c, "conversations.id=?", conversationId)

	if result.Error != nil {
		s.Logger.Println(result.Error)
		return ctx.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	resp := openapi.Conversation{
		Id: &c.ID,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetConversationMessages(ctx echo.Context, conversationId string) error {
	messages := []Message{}
	result := s.DB.Limit(20).
		Order("created_at desc").
		Where("conversation_id=?", conversationId).
		Preload("User").
		Find(&messages)
	if result.Error != nil {
		s.Logger.Println(result.Error)
		return ctx.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (s *Server) CreateUser(ctx echo.Context) error {
	requestBody := openapi.CreateUserRequestBody{}
	err := ctx.Bind(&requestBody)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	id := uuid.New().String()
	u := User{ID: id, Name: *requestBody.Name}
	result := s.DB.Create(&u)
	if result.Error != nil {
		return result.Error
	}

	return ctx.JSON(http.StatusCreated, u)
}

func (s *Server) ws(ctx echo.Context) error {
	conversationId := ctx.QueryParam("conversationId")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	conn := connection{conversationId: conversationId, ws: ws}
	if pingPongEnabled {
		if err := conn.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			s.Logger.Println(err)
		}
		conn.ws.SetPongHandler(conn.pongHandler)
	}
	s.Connections <- conn

	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.Logger.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}

		msg := Message{}
		json.Unmarshal(body, &msg)

		result := s.DB.Create(&msg)
		if result.Error != nil {
			s.Logger.Println(result.Error)
		}

		s.NewMessagesChan <- msg
	}

	s.Disconnects <- conn

	return nil
}
