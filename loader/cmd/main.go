package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/goombaio/namegenerator"
	"golang.org/x/net/websocket"
)

const (
	NumOfConversations        = 500
	NumOfUsersPerConversation = 2
)

type connectionManager struct {
	conns map[string]*websocket.Conn
	sync.RWMutex
}

type User struct {
	ID   string `json:"uuid"`
	Name string `json:"name"`
}

type Message struct {
	From           string    `json:"from"`
	CreatedAt      time.Time `json:"created_at"`
	Body           string    `json:"body"`
	ConversationId string    `json:"conversation_id"`
}

type Conversation struct {
	ID        string    `json:"uuid"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func closeConnections(connections map[string]*websocket.Conn) {
	for _, conn := range connections {
		conn.Close()
	}
}

func closeCoversations(numOfConversations int, wg *sync.WaitGroup) {
	for i := 0; i < numOfConversations; i++ {
		wg.Done()
	}
}

type loader struct {
	cm connectionManager
}

func main() {
	var wg sync.WaitGroup

	ctx := context.Background()
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// closeConnections(connections)
			closeCoversations(NumOfConversations, &wg)
		case syscall.SIGTERM:
			// closeConnections(connections)
		case syscall.SIGQUIT:
			// closeConnections(connections)
		}
	}()

	l := loader{cm: connectionManager{conns: map[string]*websocket.Conn{}}}

	for i := 0; i < NumOfConversations; i++ {
		wg.Add(1)

		var users []User
		for j := 0; j < NumOfUsersPerConversation; j++ {
			u, err := createUser(ctx, nameGenerator.Generate())
			if err != nil {
				panic(err)
			}
			users = append(users, u)
		}

		go l.startConversation(ctx, i, users)
	}

	wg.Wait()
}

func (l *loader) startConversation(ctx context.Context, conversationCount int, users []User) {
	conversation, err := createConversation(ctx)
	if err != nil {
		panic(err)
	}

	if conversationCount == 0 {
		fmt.Println(conversation.ID)
	}

	address := "api:3000"
	for _, user := range users {
		ws, err := websocket.Dial(fmt.Sprintf("ws://%s/ws?conversationId=%s", address, conversation.ID), "", fmt.Sprintf("http://%s/ws?conversationId=%s", address, conversation.ID))
		if err != nil {
			fmt.Printf("Dial failed: %s\n", err.Error())
			os.Exit(1)
		}

		l.cm.Lock()
		l.cm.conns[user.ID] = ws
		l.cm.Unlock()
	}

	for {
		user := users[rand.Intn(len(users))] // select a user at random based on the number of users in the conversation
		message := Message{From: user.ID, ConversationId: conversation.ID, Body: randSeq(rand.Intn(100))}
		conn := l.cm.conns[user.ID]
		msgJSON, _ := json.Marshal(message)
		_, err := conn.Write(msgJSON)
		if err != nil {
			fmt.Printf(err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}

func createConversation(ctx context.Context) (Conversation, error) {
	var c Conversation
	requestURL := "http://api:3000/api/v1/conversations"
	jsonBody := []byte(``)
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return c, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, errors.New("unable to create conversation")
	}

	return c, nil
}

func createUser(ctx context.Context, name string) (User, error) {
	var u User
	requestURL := "http://api:3000/api/v1/users"
	values := User{Name: name}
	jsonBody, _ := json.Marshal(values)

	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return u, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, errors.New("unable to create user")
	}

	return u, nil
}
