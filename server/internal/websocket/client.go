package websocket

import (
	"context"
	"dungeons_helper/internal/character"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
)

type Client struct {
	Conn      *websocket.Conn
	Message   chan *Command
	Id        int64  `json:"id"`
	IdLobby   int64  `json:"idLobby"`
	Nickname  string `json:"Nickname"`
	Character *character.Character
	Context   context.Context
}

type Command struct {
	Type      string               `json:"type"`
	Payload   interface{}          `json:"payload"`
	Character *character.Character `json:"character"`
}

func (c *Client) readCommand(hub *Hub, charRepo character.Repository) {
	defer func() {
		hub.LeaveRoom <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message := string(m)
		parts := strings.Fields(message)
		if len(parts) < 3 {
			sendErrorMessage(c.Conn, err.Error())
			continue
		}

		command := parts[0]
		switch command {
		case "change_hp":
			characterId := parts[1]
			charId, err := strconv.Atoi(characterId)
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			health := parts[2]
			newHealth, err := strconv.Atoi(health)
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			err = charRepo.UpdateCharacterHpById(c.Context, int64(charId), int64(newHealth))
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			c.Character, err = charRepo.GetCharacterById(c.Context, int64(charId))
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}

			cmd := &Command{
				Type:      "updateChar",
				Payload:   c.IdLobby,
				Character: c.Character,
			}
			hub.Broadcast <- cmd
			fmt.Println(c.Character)
		case "change_exp":
			characterId := parts[1]
			charId, err := strconv.Atoi(characterId)
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			exp := parts[2]
			newExp, err := strconv.Atoi(exp)
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			err = charRepo.UpdateCharacterHpById(c.Context, int64(charId), int64(newExp))
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}
			c.Character, err = charRepo.GetCharacterById(c.Context, int64(charId))
			if err != nil {
				sendErrorMessage(c.Conn, err.Error())
				continue
			}

			cmd := &Command{
				Type:      "updateChar",
				Payload:   c.IdLobby,
				Character: c.Character,
			}
			hub.Broadcast <- cmd
			fmt.Println(c.Character)
		default:
			// Обработка неизвестной команды
			fmt.Println("Unknown command:", command)
		}
	}
}
