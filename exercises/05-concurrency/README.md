# 練習 05：併發

> 對應課程：[`lessons/05-concurrency`](../../lessons/05-concurrency)　｜　檔案：`exercise.go`

**這一章的測試一定要加 `-race`：**

```bash
go test -race -v ./exercises/05-concurrency
```

資料競爭（data race）的 bug 平常可能完全看不出來，`-race` 會在執行時偵測並告訴你是哪兩行在搶同一塊記憶體。

**選擇工具的原則**

| 情境 | 工具 |
| --- | --- |
| 等一群 goroutine 做完 | `sync.WaitGroup` |
| 保護共享資料（map、計數器） | `sync.Mutex` / `sync.RWMutex` |
| 在 goroutine 之間傳遞資料、表示「完成」 | channel |
| 同時等多個事件、逾時 | `select` |
| 取消、逾時往下傳遞 | `context.Context` |

---

## 5.1 ★★ `WorkerPool(inputs []int, workers int, fn func(int) int) []int`

**規格**
- 用**固定** `workers` 個 goroutine 處理所有輸入（測試會量測同時執行數，不能超過、也不能只有 1 個）。
- 結果順序與輸入相同：`out[i] = fn(inputs[i])`。
- `workers < 1` 時視為 1。

**提示（一種做法）**
1. 建立 `jobs := make(chan int)`，傳送的是**索引**而不是值。
2. 啟動 `workers` 個 goroutine，各自 `for i := range jobs { out[i] = fn(inputs[i]) }`。
3. 主 goroutine 把所有索引送進 `jobs`，然後 `close(jobs)`。
4. `wg.Wait()` 等全部完成。

**思考題**：多個 goroutine 同時寫 `out` 這個切片，為什麼不算資料競爭？（提示：它們寫的是不同的元素。）

**延伸挑戰**：如果 `fn` 可能回傳 error，改成 `func(int) (int, error)`，在第一個錯誤發生時就停止派發新工作。去看看 `golang.org/x/sync/errgroup` 怎麼做。

---

## 5.2 ★★ `WithTimeout(fn func() int, d time.Duration) (int, error)`

**規格**：在 goroutine 中執行 `fn`；`d` 時間內完成就回傳結果，否則**立即**回傳 `ErrTimeout`。

**提示**
```go
select {
case v := <-done:
	...
case <-time.After(d):
	...
}
```

**陷阱：goroutine 洩漏**。如果 `done` 是無緩衝 channel，逾時之後 `fn` 的 goroutine 做完要送出結果時，已經沒有人在接收，它會**永遠卡住**。解法：`make(chan int, 1)`。

**學到的概念**：`select`、`time.After`、緩衝 channel、goroutine 洩漏

---

## 5.3 ★★ `FanIn(chs ...<-chan int) <-chan int`

**規格**：把多個 channel 的資料合併到一個輸出 channel。**所有**輸入都關閉後，輸出才關閉。沒有輸入時，輸出要立刻關閉。

**提示**
- 每個輸入開一個 goroutine 轉送資料。
- 再開一個 goroutine 執行 `wg.Wait(); close(out)`。
- 為什麼關閉 `out` 的動作不能放在主函式裡直接做？（主函式必須先把 `out` 回傳給呼叫端，呼叫端才能開始接收。）

**觀念**：`<-chan int` 是「只能接收」的 channel 型別；用方向性型別讓編譯器幫你檢查誤用。

---

## 5.4 ★★★ `Generate(ctx context.Context, start int) <-chan int`

**規格**：不斷送出 `start, start+1, start+2, ...`。`ctx` 被取消後，goroutine 要結束並**關閉** channel（測試會檢查 1 秒內關閉）。

**提示**
```go
select {
case out <- n:
case <-ctx.Done():
	return
}
```
送出的動作也要放在 `select` 裡——否則沒人接收時 goroutine 會卡在送出那一行，永遠看不到取消。

**學到的概念**：`context` 取消、防止 goroutine 洩漏、`defer close(out)`

**延伸挑戰**：寫 `Take(ctx, in <-chan int, n int) <-chan int` 和 `Filter(ctx, in, pred)`，把它們串成管線：`Take(ctx, Filter(ctx, Generate(ctx, 1), isPrime), 10)` 取前 10 個質數。

---

## 5.5 ★★★ `Cache[K comparable, V any]`

**規格**
- `NewCache`、`Get`、`Set`：可以被多個 goroutine 同時安全使用。
- `GetOrCompute(key, compute)`：key 存在就直接回傳；不存在就呼叫 `compute()` 並存起來。
- **關鍵要求**：100 個 goroutine 同時對同一個 key 呼叫 `GetOrCompute`，`compute` 只能被呼叫**一次**，其他人等它算完並拿到同一個結果。

**提示**
- 只用一把 `Mutex` 鎖住整個 `compute`，答案是對的，但算一個慢的 key 時，**其他 key** 也全部被卡住。能不能做得更好？
- 一種做法：map 裡存的是 `*entry`，每個 entry 有一個 `ready chan struct{}`。
  第一個人建立 entry 後放開鎖、開始計算，算完 `close(ready)`；其他人拿到 entry 後 `<-ready` 等待。
- 也可以研究 `sync.Once` 或 `golang.org/x/sync/singleflight`。

**學到的概念**：泛型結構、細粒度鎖、用 `close(channel)` 廣播通知、「快取擊穿」問題

**延伸挑戰**：為每個項目加上過期時間（TTL），以及最大容量 + LRU 淘汰。
