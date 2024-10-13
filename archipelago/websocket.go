package archipelago

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func start_websocket() (chan os.Signal, chan []byte, chan []byte, chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: os.Getenv("AP_HOST"), Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	for err != nil {
		log.Println("dial:", err)
		time.Sleep(5 * time.Second)
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	}
	log.Println("connected")

	done := make(chan struct{})

	messages := make(chan []byte)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			// log.Printf("recv: %s", message)
			messages <- message
		}
	}()

	sender := make(chan []byte)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-sender:
				err := c.WriteMessage(websocket.TextMessage, []byte(t))
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()
	return interrupt, sender, messages, done
}
