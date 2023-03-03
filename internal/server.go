package internal

import (
	"context"
	"encoding/json"
	"messanger/internal/openapi"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	Logger        *logrus.Logger
	DB            *gorm.DB
	WSConnections map[*websocket.Conn]bool
}

var (
	upgrader = websocket.Upgrader{}
)

func NewServer(ctx context.Context, logger *logrus.Logger, DB *gorm.DB) *Server {
	s := Server{
		Logger:        logger,
		DB:            DB,
		WSConnections: map[*websocket.Conn]bool{},
	}

	go func() {
		for {
			s.Logger.Printf("There is currently: %d WS connection(s)", len(s.WSConnections))
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
	s.WSConnections[ws] = true

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

			result := s.DB.Create(&msg)
			if result.Error != nil {
				s.Logger.Println(result.Error)
			}
		}()
	}

	delete(s.WSConnections, ws)

	ws.Close()

	return nil
}
