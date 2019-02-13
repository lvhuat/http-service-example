package httpserver

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		var data = struct {
			Event string `json:"event"`
			Ts    int64  `json:"ts"`
		}{
			Event: "close",
			Ts:    time.Now().UnixNano() / int64(time.Millisecond),
		}
		c.WriteJSON(gzipBytes(data))
	}()

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if mt != websocket.BinaryMessage {
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func gzipBytes(data interface{}) []byte {
	buffer := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buffer)
	defer w.Close()

	body, _ := json.Marshal(data)
	w.Write(body)
	w.Flush()
	return buffer.Bytes()
}
