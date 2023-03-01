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

type User struct {
	ID   string `json:"uuid"`
	Name string `json:"name"`
}

type Message struct {
	From           string    `json:"from_user_id"`
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

var connections map[string]*websocket.Conn

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

func main() {
	var wg sync.WaitGroup

	ctx := context.Background()
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	connections = map[string]*websocket.Conn{}
	numOfConversations := 100
	//minNumOfUsersPerConversation := 2
	//maxNumOfUsersPerConversation := 4

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			closeConnections(connections)
			closeCoversations(numOfConversations, &wg)
		case syscall.SIGTERM:
			closeConnections(connections)
		case syscall.SIGQUIT:
			closeConnections(connections)
		}
	}()

	for i := 0; i < numOfConversations; i++ {
		wg.Add(1)

		//numOfUsers := rand.Intn((maxNumOfUsersPerConversation - minNumOfUsersPerConversation)) + minNumOfUsersPerConversation
		numOfUsers := 2
		var users []User
		for j := 1; j < numOfUsers; j++ {
			u, err := createUser(ctx, nameGenerator.Generate())
			if err != nil {
				panic(err)
			}
			users = append(users, u)
		}

		go startConversation(ctx, i, users)
	}

	wg.Wait()
}

func startConversation(ctx context.Context, conversationCount int, users []User) {
	conversation, err := createConversation(ctx)
	if err != nil {
		panic(err)
	}

	address := "api:3000"
	for _, user := range users {
		ws, err := websocket.Dial(fmt.Sprintf("ws://%s/ws", address), "", fmt.Sprintf("http://%s/", address))
		if err != nil {
			fmt.Printf("Dial failed: %s\n", err.Error())
			os.Exit(1)
		}

		connections[user.ID] = ws
	}

	for {
		user := users[rand.Intn(len(users))] // select a user at random based on the number of users in the conversation
		message := Message{From: user.ID, ConversationId: conversation.ID, Body: randSeq(rand.Intn(100))}
		//fmt.Printf("(%d) %s: %s\n", conversationCount, user.Name, message.Body)
		conn := connections[user.ID]
		msgJSON, _ := json.Marshal(message)
		_, err := conn.Write(msgJSON)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("Num. of Messages: %d\n", len(conversation.Messages))
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
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
