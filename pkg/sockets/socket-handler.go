package sockets

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	pager_client "pager-services/pkg/api/pager_api/client"
	pager_common "pager-services/pkg/api/pager_api/common"
	"pager-services/pkg/client"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/transfers"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func unRegisterAndCloseConnection(c *Client) {
	c.hub.unregister <- c
	c.webSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	c.webSocketConnection.SetReadLimit(maxMessageSize)
	c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocketConnection.SetPongHandler(func(string) error { c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	type chatlistResponseStruct struct {
		Type     string      `json:"type"`
		Chatlist interface{} `json:"chatlist"`
	}

	ctx := context.Background()

	switch socketEventPayload.EventName {

	case "join":
		userID := (socketEventPayload.EventPayload).(string)
		//userID := payload.UserId
		//publicKey := payload.PublicKey
		userDetails := &pager_common.PagerProfile{}
		if err := transfers.ReadDataByID(ctx, mongo_ops.CollectionsPoll.ProfileCollection, userID, userDetails); err != nil {
			return
		}
		if userDetails == nil {
			log.Println("An invalid user with userID " + userID + " tried to connect to Chat Server.")
		} else {
			newUserOnlinePayload := SocketEventStruct{
				EventName: "user-connected",
				EventPayload: UserDetailsResponsePayloadStruct{
					Online: true,
					UserId: userDetails.UserId,
					Login:  userDetails.Login,
				},
			}
			if err := UpdateUserOnlineStatusByUserID(userID, true); err != nil {
				return
			}
			BroadcastSocketEventToAllClientExceptMe(client.hub, newUserOnlinePayload, userDetails.UserId)
		}
	case "watch-user":
		userId := (socketEventPayload.EventPayload.(map[string]interface{})["userId"]).(string)
		targetId := (socketEventPayload.EventPayload.(map[string]interface{})["targetId"]).(string)
		userDetails := &pager_common.PagerProfile{}
		if err := transfers.ReadDataByID(ctx, mongo_ops.CollectionsPoll.ProfileCollection, targetId, userDetails); err != nil {
			return
		}
		if userDetails == nil {
			log.Println("An invalid user with userID " + targetId)
		} else {
			userPayload := SocketEventStruct{
				EventName: "user-info",
				EventPayload: UserDetailsResponsePayloadStruct{
					UserId: userDetails.UserId,
					Login:  userDetails.Login,
				},
			}
			EmitToSpecificClient(client.hub, userPayload, userId)
		}
	case "disconnect":
		if socketEventPayload.EventPayload != nil {

			userID := (socketEventPayload.EventPayload).(string)
			userDetails := &pager_common.PagerProfile{}
			if err := transfers.ReadDataByID(ctx, mongo_ops.CollectionsPoll.ProfileCollection, userID, userDetails); err != nil {
				return
			}
			if err := UpdateUserOnlineStatusByUserID(userID, false); err != nil {
				return
			}

			BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
				EventName: "user-disconnected",
				EventPayload: UserDetailsResponsePayloadStruct{
					Online:         false,
					UserId:         userDetails.UserId,
					Login:          userDetails.Login,
					LastSeenMillis: time.Now().UnixMilli(),
				},
			})
		}
	}
}

func (c *Client) readPump(userId string) {
	var socketEventPayload SocketEventStruct

	// Unregistering the client and closing the connection
	defer unRegisterAndCloseConnection(c)

	// Setting up the Payload configuration
	setSocketPayloadReadConfig(c)

	for {
		// ReadMessage is a helper method for getting a reader using NextReader and reading from that reader to a buffer.
		_, payload, err := c.webSocketConnection.ReadMessage()

		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)

		if decoderErr != nil {
			if _, err := client.ChangeConnectionStateBody(context.Background(), userId, &pager_client.ConnectionRequest{
				LastStampMillis: time.Now().UnixMilli(),
				Online:          false,
			}); err != nil {
				return
			}
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}

		//  Getting the proper Payload to send the client
		handleSocketPayloadEvents(c, socketEventPayload)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.webSocketConnection.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:

			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(c.send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// CreateNewSocketUser creates a new socket user
func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, userID string) {
	client := &Client{
		hub:                 hub,
		webSocketConnection: connection,
		send:                make(chan SocketEventStruct),
		Id:                  userID,
	}

	go client.writePump()
	go client.readPump(userID)

	client.hub.register <- client
}

// HandleUserRegisterEvent will handle the Join event for New socket users
func HandleUserRegisterEvent(hub *Hub, client *Client) {
	hub.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName:    "join",
		EventPayload: client.Id,
	})
}

// HandleUserDisconnectEvent will handle the Disconnect event for socket users
func HandleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.clients[client]
	if ok {
		delete(hub.clients, client)
		close(client.send)

		handleSocketPayloadEvents(client, SocketEventStruct{
			EventName:    "disconnect",
			EventPayload: client.Id,
		})
	}
}

// EmitToSpecificClient will emit the socket event to specific socket user
func EmitToSpecificClient(hub *Hub, payload SocketEventStruct, userID string) {
	for client := range hub.clients {
		if client.Id == userID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(hub *Hub, payload SocketEventStruct) {
	for client := range hub.clients {
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}

// BroadcastSocketEventToAllClientExceptMe will emit the socket events to all socket users,
// except the user who is emitting the event
func BroadcastSocketEventToAllClientExceptMe(hub *Hub, payload SocketEventStruct, myUserID string) {
	for client := range hub.clients {
		if client.Id != myUserID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}
