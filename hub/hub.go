package hub

type hub struct {
	connections map[*conn]bool
	connect     chan *conn
	disconnect  chan *conn

	broadcast chan *message
}

var h = &hub{
	make(map[*conn]bool),
	make(chan *conn),
	make(chan *conn),
	make(chan *message),
}

func Run() {
	for {
		select {
		case conn := <-h.connect:
			h.connections[conn] = true

		case conn := <-h.disconnect:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}

		case msg := <-h.broadcast:
			for conn := range h.connections {
				select {
				case conn.send <- msg:
					// move on

				default:
					close(conn.send)
					delete(h.connections, conn)
				}
			}
		}
	}
}
