//go:build solution

// 參考解答。建議自己寫完再看！
package chat

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
)

// ---------- 10.1 Echo ----------

func ServeEcho(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil // 正常關閉
		}
		if err != nil {
			return err
		}
		go func() { // 每個連線一個 goroutine：Go 寫網路程式最自然的模式
			defer conn.Close()
			sc := bufio.NewScanner(conn)
			for sc.Scan() {
				if _, err := fmt.Fprintln(conn, sc.Text()); err != nil {
					return
				}
			}
		}()
	}
}

// ---------- 10.2 Hub ----------

type subscriber struct {
	ch   chan Message
	once sync.Once
}

type Hub struct {
	mu     sync.RWMutex
	topics map[string]map[*subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{topics: make(map[string]map[*subscriber]struct{})}
}

func (h *Hub) Subscribe(topic string, buffer int) (<-chan Message, func()) {
	sub := &subscriber{ch: make(chan Message, buffer)}
	h.mu.Lock()
	if h.topics[topic] == nil {
		h.topics[topic] = make(map[*subscriber]struct{})
	}
	h.topics[topic][sub] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		sub.once.Do(func() {
			h.mu.Lock()
			delete(h.topics[topic], sub)
			if len(h.topics[topic]) == 0 {
				delete(h.topics, topic)
			}
			// 在持有鎖的情況下關閉：Publish 也持有（讀）鎖才會送出，
			// 所以絕不會發生「對已關閉的 channel 送資料」的 panic
			close(sub.ch)
			h.mu.Unlock()
		})
	}
	return sub.ch, unsubscribe
}

func (h *Hub) Publish(topic, text string) int {
	msg := Message{Topic: topic, Text: text}
	h.mu.RLock()
	defer h.mu.RUnlock()
	delivered := 0
	for sub := range h.topics[topic] {
		select {
		case sub.ch <- msg:
			delivered++
		default: // 緩衝已滿：丟棄，不讓慢的訂閱者拖垮所有人
		}
	}
	return delivered
}

// ---------- 10.3 Chat ----------

type client struct {
	nick string
	out  chan string // 每個客戶端一個發送佇列，由專屬的 writer goroutine 寫入連線
}

type ChatServer struct {
	mu      sync.Mutex
	clients map[string]*client // 以暱稱為 key
}

func NewChatServer() *ChatServer {
	return &ChatServer{clients: make(map[string]*client)}
}

func (s *ChatServer) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

func (s *ChatServer) handle(conn net.Conn) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	if !sc.Scan() {
		return
	}
	nick := strings.TrimSpace(sc.Text())
	if nick == "" {
		fmt.Fprintln(conn, "* 暱稱不可為空")
		return
	}

	c := &client{nick: nick, out: make(chan string, 32)}
	if !s.join(c) {
		fmt.Fprintln(conn, "* 暱稱已被使用")
		return
	}

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		for msg := range c.out {
			if _, err := fmt.Fprintln(conn, msg); err != nil {
				log.Printf("寫入 %s 失敗: %v", c.nick, err)
			}
		}
	}()

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
			continue
		case line == "/quit":
			goto leave
		case line == "/who":
			c.out <- "* 線上: " + strings.Join(s.names(), ", ")
		default:
			s.broadcast(c, fmt.Sprintf("[%s] %s", nick, line))
		}
	}
leave:
	s.leave(c)
	<-writerDone // 等佇列中的訊息都寫完再關閉連線
}

func (s *ChatServer) join(c *client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, taken := s.clients[c.nick]; taken {
		return false
	}
	s.clients[c.nick] = c
	c.out <- "* 歡迎 " + c.nick
	s.broadcastLocked(c, "* "+c.nick+" 加入")
	return true
}

func (s *ChatServer) leave(c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c.nick)
	close(c.out) // 已從 map 移除，之後不會再有人送訊息給它
	s.broadcastLocked(c, "* "+c.nick+" 離開")
}

func (s *ChatServer) broadcast(from *client, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.broadcastLocked(from, msg)
}

// broadcastLocked 呼叫前必須持有 s.mu。
func (s *ChatServer) broadcastLocked(from *client, msg string) {
	for _, c := range s.clients {
		if c == from {
			continue
		}
		select {
		case c.out <- msg:
		default: // 這個客戶端太慢，丟棄訊息
		}
	}
}

func (s *ChatServer) names() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	names := make([]string, 0, len(s.clients))
	for n := range s.clients {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}

// ---------- 10.4 Server-Sent Events ----------

func EventsHandler(h *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		if topic == "" {
			http.Error(w, "缺少 topic 參數", http.StatusBadRequest)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "不支援串流", http.StatusInternalServerError)
			return
		}

		msgs, unsubscribe := h.Subscribe(topic, 16)
		defer unsubscribe() // 客戶端斷線時取消訂閱，避免洩漏

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		fmt.Fprint(w, ": connected\n\n") // 以冒號開頭的是註解，瀏覽器會忽略
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done(): // 客戶端斷線
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				for _, line := range strings.Split(msg.Text, "\n") {
					fmt.Fprintf(w, "data: %s\n", line)
				}
				fmt.Fprint(w, "\n") // 空行代表一個事件結束
				flusher.Flush()     // 不 Flush 的話資料會卡在緩衝區
			}
		}
	})
}

// ---------- 10.5 Framing ----------

func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) > MaxFrameSize {
		return ErrFrameTooLarge
	}
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if _, err := w.Write(header[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func ReadFrame(r io.Reader) ([]byte, error) {
	var header [4]byte
	// io.ReadFull：TCP 一次 Read 不保證讀滿，一定要用 ReadFull
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err // 一開始就沒資料是 io.EOF；讀到一半是 io.ErrUnexpectedEOF
	}
	n := binary.BigEndian.Uint32(header[:])
	if n > MaxFrameSize {
		return nil, ErrFrameTooLarge // 先檢查再配置記憶體，避免惡意的超大長度把記憶體吃光
	}
	payload := make([]byte, n)
	if _, err := io.ReadFull(r, payload); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return payload, nil
}
