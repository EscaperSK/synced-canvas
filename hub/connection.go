package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

type conn struct {
	*websocket.Conn

	send chan *message
}

func AddConnection(socketConn *websocket.Conn) *conn {
	c := &conn{
		socketConn,
		make(chan *message),
	}

	h.connect <- c

	return c
}

func (c *conn) Read(update func(int, bool), clear func()) {
	defer func() {
		h.disconnect <- c
		c.Close()
	}()

	for {
		var msg message

		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("Read failed:", err.Error())
			break
		}

		switch msg.Command {
		case "draw":
			update(msg.Index, msg.Fill)
		case "clear":
			clear()
		default:
			log.Println("Unknown command:", msg.Command)
			continue
		}

		h.broadcast <- &msg
	}
}

func (c *conn) Write() {
	defer func() {
		c.Close()
	}()

	for {
		msg, ok := <-c.send
		if !ok {
			c.WriteMessage(websocket.CloseMessage, nil)
			return
		}

		err := c.WriteJSON(msg)
		if err != nil {
			log.Println("Write failed:", err.Error())
			return
		}

		n := len(c.send)
		for range n {
			err := c.WriteJSON(<-c.send)
			if err != nil {
				log.Println("Write failed:", err.Error())
				continue
			}
		}
	}
}
