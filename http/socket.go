package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

var socketServer *socketio.Server

func NewSocketIO(router *gin.Engine) {
	var err error
	socketServer, err = socketio.NewServer(nil)

	if err != nil {
		panic(fmt.Sprint("Cannot create socketio ", err.Error()))
	}

	router.GET("/socket.io/*any", gin.WrapH(socketServer))
	router.POST("/socket.io/*any", gin.WrapH(socketServer))

	socketServer.OnConnect("/", onConnect)
	socketServer.OnDisconnect("/", OnDisconnect)

	socketServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

}

func GetSocketServer() *socketio.Server {
	if socketServer == nil {
		panic("SocketServer nil")
	}

	return socketServer
}

func onConnect(s socketio.Conn) error {
	s.SetContext("")
	log.Println("Client connected", s.ID())
	return nil
}

func OnDisconnect(s socketio.Conn, reason string) {
	fmt.Println("closed", reason)
}
