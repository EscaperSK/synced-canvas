package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"socket/hub"

	"github.com/gorilla/websocket"
)

var templ = template.Must(template.ParseFiles("public/page.html"))

var grid = []bool{}

const n, m = 800, 800
const dotSize = 2

func Run() {
	for i := 0; i < n*m; i++ {
		grid = append(grid, false)
	}

	go hub.Run()

	upgrader := websocket.Upgrader{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s, err := json.Marshal(grid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			N       int
			M       int
			DotSize int
			Grid    string
		}{n, m, dotSize, string(s)}

		err = templ.ExecuteTemplate(w, "page.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err.Error())
			return
		}

		c := hub.AddConnection(conn)

		updateGrid := func(index int, fill bool) {
			grid[index] = fill
		}
		clearGrid := func() {
			for i := range grid {
				grid[i] = false
			}
		}

		go c.Write()
		go c.Read(updateGrid, clearGrid)
	})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
