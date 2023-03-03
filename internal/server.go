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

type Server struct {
	Logger               *logrus.Logger
	DB                   *gorm.DB
	WSConnections        map[*websocket.Conn]bool
	MessagesChannel      chan Message
	ConnectionChannel    chan *websocket.Conn
	DisconnectionChannel chan *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{}
)

func NewServer(ctx context.Context, logger *logrus.Logger, DB *gorm.DB) *Server {
	m := make(chan Message)
	c := make(chan *websocket.Conn)
	d := make(chan *websocket.Conn)

	s := Server{
		Logger:               logger,
		DB:                   DB,
		WSConnections:        map[*websocket.Conn]bool{},
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
		for ws := range s.ConnectionChannel {
			mutex.Lock()
			s.WSConnections[ws] = true
			mutex.Unlock()
		}
	}()

	go func() {
		var mutex = &sync.Mutex{}
		for ws := range s.DisconnectionChannel {
			mutex.Lock()
			delete(s.WSConnections, ws)
			mutex.Unlock()
		}
	}()

	numOfWorkers := 75
	for w := 1; w <= numOfWorkers; w++ {
		go func(w int) {
			messageCache := []Message{}
			for msg := range s.MessagesChannel {
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
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	s.ConnectionChannel <- ws

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

	s.DisconnectionChannel <- ws
	ws.Close()

	return nil
}
