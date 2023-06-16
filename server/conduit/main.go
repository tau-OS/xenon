package conduit

import (
	"github.com/puzpuzpuz/xsync/v2"

	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
)

var conduits = xsync.NewMapOf[ConduitService]()

type websocketChannel struct {
	*websocket.Conn
}

func (c websocketChannel) Send(msg []byte) error {
	println(string(msg))
	return c.WriteMessage(websocket.TextMessage, msg)
}

func (c websocketChannel) Recv() ([]byte, error) {
	_, bytes, err := c.ReadMessage()
	println(string(bytes))
	return bytes, err
}

func (c websocketChannel) Close() error {
	return c.Conn.Close()
}

func HandleWebSocketConnection(conn *websocket.Conn) {
	user := conn.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)

	//TODO: Handle validation errors

	deviceName := conn.Query("deviceName")
	if deviceName == "" {
		conn.WriteJSON(map[string]string{})
		return
	}

	devicePublicKey := conn.Query("publicKey")
	if devicePublicKey == "" {
		conn.WriteJSON(map[string]string{})
		return
	}

	conduit, _ := conduits.LoadOrStore(userId, ConduitService{})
	server := conduit.NewRPCServer(deviceName, devicePublicKey)
	if err := server.Start(websocketChannel(websocketChannel{conn})).Wait(); err != nil {
		println(err.Error())
	}
}
