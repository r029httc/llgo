# 練習 10：網路通訊

> 檔案：`exercise.go`（你要寫的）、`types.go`（已定義的型別）
>
> ```bash
> go test -race -v ./exercises/10-network
> ```

Go 最初就是為了寫網路服務而設計的。本章從最底層的 TCP 連線開始，一路做到瀏覽器即時推送：

```
            ┌──────────────── 你的 Go 程式 ─────────────────┐
  TCP 客戶端 ──► 10.1 Echo / 10.3 聊天室 (net.Listener)      │
            │                    │                          │
            │              10.2 Hub（程序內發布/訂閱，用 channel 傳訊息）
            │                    │                          │
  瀏覽器 ◄──── 10.4 Server-Sent Events（HTTP 長連線推送）    │
            │                                               │
  其他服務 ◄──► 10.5 長度前綴封包（自訂二進位協定）          │
            └───────────────────────────────────────────────┘
```

**核心模式：每個連線一個 goroutine。** goroutine 非常便宜（初始只要幾 KB），一台機器同時處理數萬個連線是很平常的事——這是 Go 在雲端服務領域流行的主因。

**測試小技巧**
- `net.Listen("tcp", "127.0.0.1:0")`：port 0 代表讓作業系統挑一個空的 port，測試就不會互相衝突。
- `conn.SetReadDeadline(...)`：避免測試因為等不到資料而永遠卡住。
- `net.Pipe()`：在記憶體中建立一對相連的連線，不經過網路。

---

## 10.1 ★★ `ServeEcho(ln net.Listener) error` — TCP Echo 伺服器

**規格**
- 接受所有連線，把每個客戶端送來的**每一行**原封不動送回去。
- 多個客戶端可以同時連線，互不影響。
- `ln` 被關閉時回傳 `nil`（這是正常的關機流程，不是錯誤）。

**提示**
```go
for {
	conn, err := ln.Accept()
	if errors.Is(err, net.ErrClosed) { return nil }
	if err != nil { return err }
	go handle(conn) // 不能直接呼叫 handle(conn)，否則同時只能服務一個客戶端
}
```
- `handle` 裡：`defer conn.Close()`，用 `bufio.NewScanner(conn)` 逐行讀取，`fmt.Fprintln(conn, line)` 寫回。

**自己動手試**：寫一個 `cmd/echo/main.go` 呼叫 `ServeEcho`，然後用 `nc localhost 9000`（或 `telnet`）連上去打字。

---

## 10.2 ★★★ `Hub` — 發布/訂閱（pub/sub）訊息中心

這是很多真實系統的核心元件：聊天室、通知系統、即時儀表板、事件匯流排。

**規格**
| 方法 | 行為 |
| --- | --- |
| `Subscribe(topic, buffer)` | 回傳一個接收訊息的 channel（緩衝大小 `buffer`）與「取消訂閱」函式 |
| 取消訂閱函式 | 把訂閱者移除並**關閉** channel；呼叫多次也不能 panic |
| `Publish(topic, text)` | 送給該 topic 的所有訂閱者，回傳成功送達的數量 |

**三個測試會檢查的陷阱**
1. **慢訂閱者**：某個訂閱者不讀取、緩衝已滿時，`Publish` **不能卡住**（否則一個慢的客戶端會拖垮所有人）。要丟棄給它的訊息：
   ```go
   select {
   case ch <- msg:  // 送得出去
   default:         // 緩衝滿了：丟棄
   }
   ```
2. **對已關閉的 channel 送資料會 panic**：如果「取消訂閱（close）」和「Publish（send）」同時發生，就可能出事。解法：兩者都在同一把鎖的保護下進行（Publish 用讀鎖、取消訂閱用寫鎖）。
3. **重複取消訂閱**：第二次 `close` 同一個 channel 也會 panic。用 `sync.Once` 包起來。

**提示**：資料結構可以是 `map[string]map[*subscriber]struct{}`（topic → 訂閱者集合），用 `sync.RWMutex` 保護。

**延伸挑戰**
- 統計每個訂閱者被丟棄了幾則訊息，超過門檻就自動踢掉。
- 支援萬用字元訂閱：`Subscribe("sports.*")`。
- 研究真正的訊息佇列系統：NATS（Go 寫的）、Redis Pub/Sub、Kafka。它們和你的 Hub 差在哪裡？（提示：跨機器、持久化、送達保證。）

---

## 10.3 ★★★ `ChatServer` — TCP 多人聊天室

**協定**（每則訊息都是一行文字）

| 時機 | 誰收到 | 內容 |
| --- | --- | --- |
| 連線後送出的**第一行** | — | 當作暱稱 |
| 暱稱為空（去掉空白後） | 自己 | `* 暱稱不可為空`，然後斷線 |
| 暱稱已被使用 | 自己 | `* 暱稱已被使用`，然後斷線 |
| 加入成功 | 自己 | `* 歡迎 <暱稱>` |
| 加入成功 | 其他人 | `* <暱稱> 加入` |
| 送出一般訊息 | 其他人（不含自己） | `[<暱稱>] <內容>` |
| 送出 `/who` | 自己 | `* 線上: a, b, c`（依字母排序，以 `, ` 分隔） |
| 送出 `/quit` 或直接斷線 | 其他人 | `* <暱稱> 離開` |

空白行忽略。

**建議架構**
```
          ┌── reader goroutine（handle）：讀取這個客戶端的輸入 → 廣播給其他人
連線 ─────┤
          └── writer goroutine：從 out channel 取訊息 → 寫入連線

ChatServer { mu sync.Mutex; clients map[string]*client }
client     { nick string; out chan string }
```
- 每個客戶端有自己的 `out` channel（有緩衝）與專屬的 writer goroutine。**廣播時只把訊息丟進 channel**，不直接寫入連線——這樣一個網路很慢的客戶端不會卡住廣播。
- 和 10.2 一樣，廣播用 `select` + `default` 做非阻塞送出。
- 離開時：先從 map 移除（持有鎖）→ 再 `close(out)` → 等 writer 結束 → 關閉連線。順序錯了就可能對已關閉的 channel 送資料。

**自己動手試**：寫 `cmd/chat/main.go` 啟動伺服器，開三個終端機都 `nc localhost 9000`，就是一個真的聊天室了！

**延伸挑戰**
- 新增 `/nick 新名字`、`/msg 對象 私訊內容`。
- 聊天室（rooms）：`/join #golang`，訊息只廣播給同一個房間——可以改用 10.2 的 Hub 實作。
- 閒置 5 分鐘自動斷線（`conn.SetReadDeadline`）。
- 優雅關機：`Shutdown(ctx)` 通知所有人「伺服器即將關閉」後再斷線。
- 把訊息存進練習 09 的資料庫，新加入的人可以看到最近 20 則歷史訊息。

---

## 10.4 ★★ `EventsHandler(h *Hub) http.Handler` — Server-Sent Events

SSE 是瀏覽器原生支援的**伺服器推送**技術：一個不會結束的 HTTP 回應，伺服器有新資料就寫一段過去。股價、通知、AI 聊天的逐字輸出，常常都是用它做的。比 WebSocket 簡單，而且只需要標準函式庫。

**規格**（`GET /events?topic=news`）
- 沒有 `topic` 參數 → 400。
- 回應標頭 `Content-Type: text/event-stream`。
- 訂閱成功後，先送出 `: connected\n\n` 並 **Flush**（冒號開頭是註解，瀏覽器會忽略；測試用它確認已經訂閱好了）。
- 每則訊息送出 `data: <內容>\n\n`。訊息含換行時，**每一行**都要有自己的 `data: ` 前綴：
  ```
  data: 第一行
  data: 第二行

  ```
- 客戶端斷線（`r.Context().Done()`）時結束 handler，**並取消訂閱**（測試會檢查有沒有洩漏）。

**提示**
- `flusher, ok := w.(http.Flusher)`：每寫完一個事件就 `flusher.Flush()`，否則資料會卡在緩衝區裡，瀏覽器收不到。
- 主迴圈：
  ```go
  for {
  	select {
  	case <-r.Context().Done():
  		return
  	case msg := <-msgs:
  		// 寫入並 Flush
  	}
  }
  ```

**在瀏覽器裡試試**：前端只需要這幾行：
```html
<ul id="log"></ul>
<script>
  const es = new EventSource("/events?topic=news");
  es.onmessage = (e) => {
    const li = document.createElement("li");
    li.textContent = e.data;
    document.getElementById("log").append(li);
  };
</script>
```
再加一個 `POST /publish` handler 呼叫 `hub.Publish`，用 `curl -d '有新消息' localhost:8080/publish` 發送，瀏覽器就會即時出現訊息。

**延伸挑戰**
- 支援 `id:` 欄位與 `Last-Event-ID` 標頭，讓斷線重連的瀏覽器補收漏掉的訊息。
- 改用 WebSocket（[`github.com/coder/websocket`](https://github.com/coder/websocket)）做雙向通訊，把 10.3 的聊天室搬到網頁上。

---

## 10.5 ★★ `WriteFrame` / `ReadFrame` — 長度前綴封包

TCP 是一條**位元組串流**，沒有「訊息」的邊界：你送了兩次 `Write("hello")`，對方可能一次讀到 `"hellohello"`，也可能分三次讀到 `"he"`、`"llohe"`、`"llo"`。
以「換行」分隔（10.1、10.3）只適用於文字；傳送任意二進位資料時，最常見的做法是在前面加上長度：

```
┌───────────────┬──────────────────────┐
│ 長度 (4 bytes, │      內容 (N bytes)   │
│  big-endian)  │                      │
└───────────────┴──────────────────────┘
例："hello" → 00 00 00 05 68 65 6c 6c 6f
```

gRPC、PostgreSQL、Redis、Kafka 等協定都使用類似的設計。

**規格**
- `WriteFrame`：內容超過 `MaxFrameSize` → `ErrFrameTooLarge`。
- `ReadFrame`：
  - 一開始就沒有資料 → `io.EOF`（代表對方正常結束）。
  - 標頭或內容讀到一半就斷了 → `io.ErrUnexpectedEOF`。
  - 標頭宣稱的長度超過 `MaxFrameSize` → `ErrFrameTooLarge`，**而且不能先配置記憶體**。

**提示**
- `binary.BigEndian.PutUint32(buf, n)` / `binary.BigEndian.Uint32(buf)`。
- **一定要用 `io.ReadFull`**，不能只呼叫一次 `r.Read`——`Read` 讀到的資料量可能比你要的少。`io.ReadFull` 本身就會回傳正確的 `io.EOF` / `io.ErrUnexpectedEOF`。

**安全性觀念**：如果先 `make([]byte, n)` 才檢查長度，攻擊者只要送 4 個 bytes `FF FF FF FF`，你的伺服器就會嘗試配置 4GB 記憶體然後當掉。**永遠不要信任從網路讀到的長度。**

**延伸挑戰**：在封包內容中放 JSON，做一個簡單的 RPC：客戶端送 `{"method":"add","params":[1,2]}`，伺服器回 `{"result":3}`。再看看標準函式庫的 `net/rpc` 怎麼做。
