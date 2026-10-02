package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// listen 在本機隨機 port 開一個 TCP listener（port 0 = 讓系統選）。
func listen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// requireServe 先用「已關閉的 listener」同步呼叫一次 serve：
// 正確的實作應立即回傳 nil；尚未實作則在這裡 panic → SKIP（而不是在 goroutine 裡讓整個測試程式當掉）。
func requireServe(t *testing.T, serve func(net.Listener) error) {
	t.Helper()
	ln := listen(t)
	ln.Close()
	if err := serve(ln); err != nil {
		t.Fatalf("listener 已關閉時應回傳 nil，得到 %v", err)
	}
}

type conn struct {
	t *testing.T
	net.Conn
	r *bufio.Reader
}

func dial(t *testing.T, addr string) *conn {
	t.Helper()
	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return &conn{t, c, bufio.NewReader(c)}
}

func (c *conn) send(line string) {
	c.t.Helper()
	if _, err := fmt.Fprintln(c, line); err != nil {
		c.t.Fatal(err)
	}
}

func (c *conn) expect(want string) {
	c.t.Helper()
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	got, err := c.r.ReadString('\n')
	if err != nil {
		c.t.Fatalf("等待 %q 時發生錯誤: %v", want, err)
	}
	if got = strings.TrimRight(got, "\r\n"); got != want {
		c.t.Fatalf("收到 %q, want %q", got, want)
	}
}

// ---------- 10.1 ----------

func TestServeEcho(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	requireServe(t, ServeEcho)

	ln := listen(t)
	done := make(chan error, 1)
	go func() { done <- ServeEcho(ln) }()

	a, b := dial(t, ln.Addr().String()), dial(t, ln.Addr().String())
	a.send("hello")
	b.send("你好")
	a.send("second line")
	a.expect("hello")
	b.expect("你好")
	a.expect("second line")

	ln.Close()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("關閉 listener 後應回傳 nil，得到 %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("關閉 listener 後 ServeEcho 沒有返回")
	}
}

// ---------- 10.2 ----------

func recv(t *testing.T, ch <-chan Message) Message {
	t.Helper()
	select {
	case m := <-ch:
		return m
	case <-time.After(time.Second):
		t.Fatal("1 秒內沒有收到訊息")
		return Message{}
	}
}

func TestHub(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewHub()
	go1, unsub1 := h.Subscribe("go", 4)
	go2, unsub2 := h.Subscribe("go", 4)
	rust, unsubRust := h.Subscribe("rust", 4)
	defer unsub2()
	defer unsubRust()

	if n := h.Publish("go", "Go 1.24 發布"); n != 2 {
		t.Errorf("Publish(go) 送達 %d 個，want 2", n)
	}
	if m := recv(t, go1); m != (Message{"go", "Go 1.24 發布"}) {
		t.Errorf("go1 收到 %+v", m)
	}
	recv(t, go2)
	select {
	case m := <-rust:
		t.Errorf("rust 訂閱者不應收到 go 的訊息：%+v", m)
	default:
	}
	if n := h.Publish("nobody", "x"); n != 0 {
		t.Errorf("沒人訂閱的 topic 應送達 0，得到 %d", n)
	}

	unsub1()
	unsub1() // 呼叫兩次也不能 panic
	if _, ok := <-go1; ok {
		t.Error("取消訂閱後 channel 應被關閉")
	}
	if n := h.Publish("go", "after"); n != 1 {
		t.Errorf("取消訂閱後應只剩 1 個訂閱者，送達 %d", n)
	}
}

func TestHubSlowSubscriber(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewHub()
	slow, unsubSlow := h.Subscribe("t", 1) // 緩衝只有 1，而且我們不讀它
	fast, unsubFast := h.Subscribe("t", 100)
	defer unsubSlow()
	defer unsubFast()

	done := make(chan struct{})
	go func() {
		for i := range 10 {
			h.Publish("t", fmt.Sprint(i))
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Publish 被慢的訂閱者卡住了（要用 select + default 做非阻塞送出）")
	}
	if len(fast) != 10 {
		t.Errorf("快的訂閱者應收到全部 10 則，得到 %d", len(fast))
	}
	if m := <-slow; m.Text != "0" {
		t.Errorf("慢的訂閱者應保留第一則，得到 %q", m.Text)
	}
}

func TestHubConcurrent(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewHub()
	_, u := h.Subscribe("warmup", 1)
	u()
	// 同時發布與訂閱/取消訂閱：用 -race 檢查資料競爭，也檢查「對已關閉 channel 送資料」的 panic
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for range 100 {
				h.Publish("t", "x")
			}
		}()
		go func() {
			defer wg.Done()
			ch, unsub := h.Subscribe("t", i%3)
			go func() {
				for range ch {
				}
			}()
			time.Sleep(time.Millisecond)
			unsub()
		}()
	}
	wg.Wait()
}

// ---------- 10.3 ----------

func startChat(t *testing.T) string {
	t.Helper()
	s := NewChatServer()
	requireServe(t, s.Serve)
	ln := listen(t)
	t.Cleanup(func() { ln.Close() })
	go s.Serve(ln)
	return ln.Addr().String()
}

func join(t *testing.T, addr, nick string) *conn {
	t.Helper()
	c := dial(t, addr)
	c.send(nick)
	c.expect("* 歡迎 " + nick)
	return c
}

func TestChatServer(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	addr := startChat(t)

	alice := join(t, addr, "alice")
	bob := join(t, addr, "bob")
	alice.expect("* bob 加入")

	alice.send("嗨 bob")
	bob.expect("[alice] 嗨 bob")

	carol := join(t, addr, "carol")
	alice.expect("* carol 加入")
	bob.expect("* carol 加入")

	carol.send("/who")
	carol.expect("* 線上: alice, bob, carol")

	bob.send("大家好")
	alice.expect("[bob] 大家好")
	carol.expect("[bob] 大家好")

	bob.send("/quit")
	alice.expect("* bob 離開")
	carol.expect("* bob 離開")

	carol.Close() // 直接斷線（沒有 /quit）也要正確處理
	alice.expect("* carol 離開")
	alice.send("/who")
	alice.expect("* 線上: alice")
}

func TestChatServerNicknames(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	addr := startChat(t)
	join(t, addr, "alice")

	dup := dial(t, addr)
	dup.send("alice")
	dup.expect("* 暱稱已被使用")

	empty := dial(t, addr)
	empty.send("   ")
	empty.expect("* 暱稱不可為空")
}

// ---------- 10.4 ----------

func TestEventsHandler(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	hub := NewHub()
	srv := httptest.NewServer(EventsHandler(hub))
	defer srv.Close()

	if resp, err := http.Get(srv.URL); err != nil || resp.StatusCode != http.StatusBadRequest {
		t.Errorf("缺少 topic 應回 400，得到 %v %v", resp.StatusCode, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", srv.URL+"/events?topic=news", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	r := bufio.NewReader(resp.Body)
	readLine := func() string {
		t.Helper()
		line, err := r.ReadString('\n')
		if err != nil {
			t.Fatalf("讀取事件串流: %v", err)
		}
		return strings.TrimRight(line, "\n")
	}
	if l := readLine(); l != ": connected" {
		t.Fatalf("第一行應為 \": connected\"，得到 %q（記得 Flush）", l)
	}
	readLine() // 空行

	hub.Publish("news", "今日頭條")
	if l := readLine(); l != "data: 今日頭條" {
		t.Errorf("收到 %q, want \"data: 今日頭條\"", l)
	}
	readLine()

	hub.Publish("news", "第一行\n第二行")
	got := []string{readLine(), readLine(), readLine()}
	if want := []string{"data: 第一行", "data: 第二行", ""}; strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("多行訊息 = %q, want %q", got, want)
	}

	cancel() // 瀏覽器關閉分頁
	deadline := time.Now().Add(2 * time.Second)
	for hub.Publish("news", "ping") != 0 {
		if time.Now().After(deadline) {
			t.Fatal("客戶端斷線後仍有訂閱者：handler 結束時要取消訂閱（defer unsubscribe()）")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// ---------- 10.5 ----------

func TestFrames(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	var buf bytes.Buffer
	msgs := [][]byte{[]byte("hello"), {}, []byte("中文訊息"), bytes.Repeat([]byte{0xff}, 70000)}
	for _, m := range msgs {
		if err := WriteFrame(&buf, m); err != nil {
			t.Fatal(err)
		}
	}
	if raw := buf.Bytes(); !bytes.Equal(raw[:9], []byte{0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o'}) {
		t.Errorf("封包格式錯誤，前 9 bytes = %v", raw[:9])
	}
	for i, want := range msgs {
		got, err := ReadFrame(&buf)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("第 %d 個封包 = %d bytes, %v", i, len(got), err)
		}
	}
	if _, err := ReadFrame(&buf); err != io.EOF {
		t.Errorf("沒有資料時應回傳 io.EOF，得到 %v", err)
	}
}

func TestFramesErrors(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if err := WriteFrame(io.Discard, make([]byte, MaxFrameSize+1)); !errors.Is(err, ErrFrameTooLarge) {
		t.Errorf("寫入超大封包應 ErrFrameTooLarge，得到 %v", err)
	}

	huge := binary.BigEndian.AppendUint32(nil, 0xFFFFFFFF) // 惡意宣稱 4GB
	if _, err := ReadFrame(bytes.NewReader(huge)); !errors.Is(err, ErrFrameTooLarge) {
		t.Errorf("讀到超大長度應 ErrFrameTooLarge（且不能真的配置 4GB），得到 %v", err)
	}

	truncated := append(binary.BigEndian.AppendUint32(nil, 10), "abc"...)
	if _, err := ReadFrame(bytes.NewReader(truncated)); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("內容不完整應 io.ErrUnexpectedEOF，得到 %v", err)
	}
	if _, err := ReadFrame(bytes.NewReader([]byte{0, 0})); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("標頭不完整應 io.ErrUnexpectedEOF，得到 %v", err)
	}
}

func TestFramesOverNetwork(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	WriteFrame(io.Discard, nil) // 尚未實作時在這裡 SKIP

	// net.Pipe 建立一對在記憶體中相連的連線，非常適合測試網路協定
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go func() { // 伺服器：讀一個封包，回傳反轉後的內容
		for {
			req, err := ReadFrame(server)
			if err != nil {
				return
			}
			for i, j := 0, len(req)-1; i < j; i, j = i+1, j-1 {
				req[i], req[j] = req[j], req[i]
			}
			WriteFrame(server, req)
		}
	}()

	for _, msg := range []string{"ping", "abcdef"} {
		if err := WriteFrame(client, []byte(msg)); err != nil {
			t.Fatal(err)
		}
		resp, err := ReadFrame(client)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte(msg)
		for i, j := 0, len(want)-1; i < j; i, j = i+1, j-1 {
			want[i], want[j] = want[j], want[i]
		}
		if !bytes.Equal(resp, want) {
			t.Errorf("回應 %q, want %q", resp, want)
		}
	}
}
