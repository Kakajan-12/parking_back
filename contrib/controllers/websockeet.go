package controllers

import (
 

	"github.com/gofiber/websocket/v2"

	"backend/contrib/models"
)
var clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan models.CarSession)

var Refresh = make(chan struct{})

func CarSessionWebsocket(c *websocket.Conn) {
	defer func() {
		delete(clients, c)
		c.Close()
	}()

	clients[c] = true

	for {
		var msg interface{}
		if err := c.ReadJSON(&msg); err != nil {
			break
		}

		if _, ok := msg.(string); ok && msg == "refresh" {
			Refresh <- struct{}{}
		} else if carSession, ok := msg.(models.CarSession); ok {

			Broadcast <- carSession
		}
	}
}

func HandleMessages() {
	for {
		select {
		case car := <-Broadcast:
			for client := range clients {
				if err := client.WriteJSON(car); err != nil {
					client.Close()
					delete(clients, client)
				}
			}

		case <-Refresh:
			for client := range clients {
				if err := client.WriteJSON("refresh"); err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}