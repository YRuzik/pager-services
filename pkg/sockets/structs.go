package sockets

import "github.com/gorilla/websocket"

// UserDetailsStruct is a universal struct for mapping the user details
type UserDetailsStruct struct {
	ID       string `bson:"_id,omitempty"`
	Login    string
	Online   bool
	SocketID string
}

// UserDetailsRequestPayloadStruct represents payload for Login and Registration request
type UserDetailsRequestPayloadStruct struct {
	Username string
	Password string
}

// UserDetailsResponsePayloadStruct represents payload for Login and Registration response
type UserDetailsResponsePayloadStruct struct {
	Login          string `json:"login"`
	UserId         string `json:"user_id"`
	Online         bool   `json:"online"`
	LastSeenMillis int64  `json:"lastSeenMillis"`
}

// SocketEventStruct struct of socket events
type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

type SocketUserConnectionPayload struct {
	UserId    string `json:"user_id"`
	PublicKey string `json:"public_key"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	send                chan SocketEventStruct
	Id                  string
}
