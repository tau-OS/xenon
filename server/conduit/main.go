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

	conduit, _ := conduits.LoadOrStore(userId, ConduitService{})
	if err := conduit.NewServer().Start(websocketChannel(websocketChannel{conn})).Wait(); err != nil {
		panic(err.Error())
	}
}
