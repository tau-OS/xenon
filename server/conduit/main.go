package conduit

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/puzpuzpuz/xsync/v2"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	closeGracePeriod = 10 * time.Second
)

var validate = validator.New()
var conduits = xsync.NewMapOf[*ConduitService]()

type websocketChannel struct {
	*websocket.Conn
}

func (c websocketChannel) Send(msg []byte) error {
	return c.WriteMessage(websocket.TextMessage, msg)
}

func (c websocketChannel) Recv() ([]byte, error) {
	messageType, bytes, err := c.ReadMessage()

	if messageType == websocket.CloseMessage {
		return nil, errors.New("connection closed")
	}

	if messageType != websocket.TextMessage {
		return nil, errors.New("invalid message type")
	}

	return bytes, err
}

func (c websocketChannel) Close() error {
	return c.Conn.Close()
}

type connectionParams struct {
	DeviceName      string `validate:"required"`
	DevicePublicKey string `validate:"required"`
}

// HandleConduitRequest handles a request to the conduit service, before the websocket connection is established.
func HandleConduitRequest(c *fiber.Ctx) error {
	params := &connectionParams{
		DeviceName:      c.Query("deviceName"),
		DevicePublicKey: c.Query("publicKey"),
	}

	if errs := validate.StructCtx(c.Context(), params); errs != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	c.Locals("connectionParams", params)

	return c.Next()
}

func HandleWebSocketConnection(conn *websocket.Conn) {
	user := conn.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)

	params := conn.Locals("connectionParams").(*connectionParams)

	conduit, _ := conduits.LoadOrStore(userId, &ConduitService{})
	server, err := conduit.NewRPCServer(params.DeviceName, params.DevicePublicKey)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error()))
		time.Sleep(closeGracePeriod)
		conn.Close()
		return
	}

	defer conduit.RemoveDevice(params.DevicePublicKey)

	conn.SetCloseHandler(func(code int, text string) error {
		server.Stop()
		return nil
	})

	if err := server.Start(websocketChannel(websocketChannel{conn})).Wait(); err != nil {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error()))
		time.Sleep(closeGracePeriod)
		conn.Close()
	}
}
