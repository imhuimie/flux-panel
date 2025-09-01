package gost_api

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type GostClient struct {
	conn    *websocket.Conn
	mu      sync.Mutex
	pending map[string]chan []byte
	url     string
}

func NewGostClient(url string) *GostClient {
	return &GostClient{
		pending: make(map[string]chan []byte),
		url:     url,
	}
}

func (c *GostClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop()
	return nil
}

func (c *GostClient) readLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			c.reconnect()
			return
		}

		var response map[string]interface{}
		if err := json.Unmarshal(message, &response); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		if id, ok := response["id"].(string); ok {
			c.mu.Lock()
			if ch, ok := c.pending[id]; ok {
				delete(c.pending, id)
				ch <- message
			}
			c.mu.Unlock()
		}
	}
}

func (c *GostClient) reconnect() {
	c.conn.Close()
	for {
		log.Println("Reconnecting to", c.url)
		if err := c.Connect(); err == nil {
			log.Println("Reconnected to", c.url)
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func (c *GostClient) Send(id string, method string, params interface{}) ([]byte, error) {
	req := map[string]interface{}{
		"id":     id,
		"method": method,
		"params": params,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)
	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(10 * time.Second):
		return nil, log.Output(2, "request timeout")
	}
}
