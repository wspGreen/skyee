package websocket_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// 升级传入的 HTTP 连接 ，并返回一个指向 WebSocket 连接的指针
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	// err = ws.WriteMessage(1, []byte("Hi Client!"))

	// if err != nil {
	// 	log.Println(err)
	// }

	gonode.Go(func() {
		reader(ws)

	})
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println("Request :", string(p))

		req := fmt.Sprintf("Response : %s", string(p))

		if err := conn.WriteMessage(messageType, []byte(req)); err != nil {
			log.Println(err)
			return
		}

	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func TestWS(t *testing.T) {
	fmt.Println("Start WS")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8081", nil))
}
