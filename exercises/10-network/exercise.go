//go:build !solution

// 練習 10：網路通訊（TCP、發布/訂閱、聊天室、Server-Sent Events、封包協定）。
// 題目說明請看 README.md。寫完務必用 go test -race 執行！
package chat

import (
	"io"
	"net"
	"net/http"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 10.1 ★★ ServeEcho 接受連線，把每個客戶端送來的每一行原封不動送回去。
// listener 被關閉時回傳 nil。
func ServeEcho(ln net.Listener) error {
	panic(todo.NotImplemented)
}

// 10.2 ★★★ Hub 是程序內的發布/訂閱（pub/sub）訊息中心。
type Hub struct {
	// TODO: 加入你需要的欄位
}

func NewHub() *Hub { panic(todo.NotImplemented) }

// Subscribe 訂閱 topic，回傳接收訊息的 channel（緩衝大小為 buffer）與取消訂閱函式。
// 取消訂閱後 channel 會被關閉；取消函式可以安全地呼叫多次。
func (h *Hub) Subscribe(topic string, buffer int) (<-chan Message, func()) {
	panic(todo.NotImplemented)
}

// Publish 把訊息送給 topic 的所有訂閱者，回傳成功送達的數量。
// 絕不能因為某個訂閱者太慢（緩衝已滿）而阻塞：直接丟棄給它的這則訊息。
func (h *Hub) Publish(topic, text string) int {
	panic(todo.NotImplemented)
}

// 10.3 ★★★ ChatServer 是 TCP 多人聊天室。協定請看 README.md。
type ChatServer struct {
	// TODO: 加入你需要的欄位
}

func NewChatServer() *ChatServer                  { panic(todo.NotImplemented) }
func (s *ChatServer) Serve(ln net.Listener) error { panic(todo.NotImplemented) }

// 10.4 ★★ EventsHandler 以 Server-Sent Events 把 Hub 的訊息即時推送給瀏覽器。
// 用法：GET /events?topic=news
func EventsHandler(h *Hub) http.Handler {
	panic(todo.NotImplemented)
}

// 10.5 ★★ WriteFrame / ReadFrame：長度前綴（length-prefixed）的二進位封包協定。
// 格式：4 bytes 大端序（big-endian）長度 + 內容。
func WriteFrame(w io.Writer, payload []byte) error {
	panic(todo.NotImplemented)
}

func ReadFrame(r io.Reader) ([]byte, error) {
	panic(todo.NotImplemented)
}
