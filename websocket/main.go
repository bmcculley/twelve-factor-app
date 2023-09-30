package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	pollPeriod = 10 * time.Second
)

var (
	origins   []string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			var origin = r.Header.Get("origin")
			for _, allowOrigin := range origins {
				if origin == allowOrigin {
					return true
				}
			}
			return false
		},
	}
)

type Message struct {
	Visitors string `json:"visitors"`
}

type Config struct {
	Addr    string
	Origins []string
}

func parseConfig(configStr string) (Config, error) {
	var jStruct Config
	err := json.Unmarshal([]byte(configStr), &jStruct)
	if err != nil {
		return Config{}, err
	}
	return jStruct, nil
}

func getVisitors() ([]uint8, error) {
	var n []uint8
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	visitorCount, err := rdb.Get(ctx, "visitorCount").Result()
	if err == redis.Nil {
		return n, errors.New("visitorCount does not exist")
	} else if err != nil {
		return n, err
	}
	s := &Message{Visitors: visitorCount}
	jVisitorCount, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	return jVisitorCount, nil
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	visitorTicker := time.NewTicker(pollPeriod)
	defer func() {
		pingTicker.Stop()
		visitorTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-visitorTicker.C:

			visitorCount, err := getVisitors()
			if err != nil {
				panic(err)
			}

			if visitorCount != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, visitorCount); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection from: %s\n", r.Header.Get("origin"))
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	visitorCount, err := getVisitors()
	if err != nil {
		panic(err)
	}

	if visitorCount != nil {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, visitorCount); err != nil {
			log.Println(err)
		}
	}
	go writer(ws)
	reader(ws)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	configStr := getEnv("WS_CONFIG", `{
		"addr": ":3000",
		"origins": [
		  "http://127.0.0.1:3000",
		  "http://localhost:3000",
		  "http://demo.ex.net:8000"
		]
	}`)
	jd, err := parseConfig(configStr)
	if err != nil {
		panic(err)
	}
	origins = jd.Origins
	http.HandleFunc("/ws", serveWs)
	server := &http.Server{
		Addr:              jd.Addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}