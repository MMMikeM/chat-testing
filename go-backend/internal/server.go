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
	conversationID string
	ws             *websocket.Conn
}

type Server struct {
	Logger               *logrus.Logger
	DB                   *gorm.DB
	WSConnections        map[string]map[*websocket.Conn]bool
	MessagesChannel      chan Message
	ConnectionChannel    chan connection
	DisconnectionChannel chan connection
}

var (
	upgrader = websocket.Upgrader{}
)

func NewServer(ctx context.Context, logger *logrus.Logger, DB *gorm.DB) *Server {
	m := make(chan Message)
	c := make(chan connection)
	d := make(chan connection)

	s := Server{
		Logger:               logger,
		DB:                   DB,
		WSConnections:        map[string]map[*websocket.Conn]bool{},
		MessagesChannel:      m,
		ConnectionChannel:    c,
		DisconnectionChannel: d,
	}

	go func() {
		for {
			s.Logger.Printf("There is currently: %d WS connection(s)", len(s.WSConnections))
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	go func() {
		var mutex = &sync.Mutex{}
		for connection := range s.ConnectionChannel {
			go func() {
				mutex.Lock()
				if s.WSConnections[connection.conversationID] == nil {
					s.WSConnections[connection.conversationID] = map[*websocket.Conn]bool{}
				}
				s.WSConnections[connection.conversationID][connection.ws] = true
				mutex.Unlock()
			}()
		}
	}()

	go func() {
		var mutex = &sync.Mutex{}
		for connection := range s.DisconnectionChannel {
			mutex.Lock()
			delete(s.WSConnections[connection.conversationID], connection.ws)
			mutex.Unlock()
		}
	}()

	numOfWorkers := 75
	for w := 1; w <= numOfWorkers; w++ {
		go func(w int) {
			messageCache := []Message{}

			for msg := range s.MessagesChannel {
				msg.ID = uuid.New().String()
				msgString, _ := json.Marshal(msg)
				for ws := range s.WSConnections[msg.ConversationId] {
					err := ws.WriteMessage(websocket.TextMessage, msgString)
					if err != nil {
						s.Logger.Println(err.Error())
					}
				}
				messageCache = append(messageCache, msg)

				if len(messageCache) > 500 {
					result := s.DB.Create(&messageCache)
					if result.Error != nil {
						s.Logger.Println(result.Error)
					}
					messageCache = []Message{}
				}
			}
		}(w)
	}

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
	conversationID := ctx.QueryParam("conversation_id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	connection := connection{conversationID: conversationID, ws: ws}
	s.ConnectionChannel <- connection

	// manage a count of connections without disturbing connection performance
	// var mutex = &sync.Mutex{}
	// go func() {
	// 	mutex.Lock()
	// 	s.WSConnections[ws] = true
	// 	mutex.Unlock()
	// }()

	for {
		// Read
		mt, body, err := ws.ReadMessage()
		if err != nil {
			s.Logger.Println(err.Error())
		}

		if err != nil || mt == websocket.CloseMessage {
			break // Exit the loop if the client tries to close the connection or the connection with the interrupted client
		}

		// separate go function for processing messages and unlocking channel
		go func() {
			msg := Message{}
			json.Unmarshal(body, &msg)

			s.MessagesChannel <- msg
		}()
	}

	// go func() {
	// 	mutex.Lock()
	// 	delete(s.WSConnections, ws)
	// 	mutex.Unlock()
	// }()

	s.DisconnectionChannel <- connection
	ws.Close()

	return nil
}
