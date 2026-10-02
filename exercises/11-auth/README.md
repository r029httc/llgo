# 練習 11：認證與安全

> 檔案：`exercise.go`（你要寫的）、`types.go`（已定義的型別與錯誤）

任何有「使用者」的服務都必須處理這些問題：密碼怎麼存？登入後怎麼證明「我是我」？怎麼擋住暴力破解？
本章做出一套完整、正確的認證元件。**安全相關的程式碼，「看起來能動」和「正確」是兩回事**——測試會實際嘗試攻擊你的實作。

```
登入流程：
  POST /login {email, password}
     │  從資料庫（練習 09）查出使用者的密碼雜湊
     ▼
  11.1 CheckPassword(hash, password) ── 失敗 → 401
     │  成功
     ▼
  11.3 Sign(Claims{UserID, Role, ExpiresAt}) → token 回傳給客戶端

之後的每個請求：
  GET /admin   Authorization: Bearer <token>
     │
  11.6 RateLimit ─────── 太頻繁 → 429
     │
  11.5 Authenticate ──── token 無效/過期 → 401   （把 Claims 放進 context）
     │
  11.5 RequireRole("admin") ── 角色不符 → 403
     │
  你的 handler：ClaimsFrom(r.Context()) 知道是誰在呼叫
```

---

## 11.1 ★ `HashPassword` / `CheckPassword` — 密碼雜湊

**規格**
- 少於 8 個**字元**（不是 bytes；8 個中文字要合法）→ `ErrWeakPassword`。
- 使用 `golang.org/x/crypto/bcrypt`，成本用套件變數 `BcryptCost`。
- 超過 72 bytes 時 bcrypt 會回傳錯誤，要把它傳回去。
- `CheckPassword` 對錯誤密碼、無效的雜湊字串都回傳 `false`。

**提示**
- `bcrypt.GenerateFromPassword([]byte(pw), BcryptCost)`、`bcrypt.CompareHashAndPassword(hash, pw)`。
- `utf8.RuneCountInString(s)` 計算字元數。

**觀念：為什麼不能用 SHA-256 存密碼？**
- **絕對不能存明文**：資料庫外洩時所有密碼直接曝光。
- **一般雜湊（MD5、SHA-256）太快了**：一張顯示卡每秒可以算數十億次，常見密碼幾秒就被暴力破解。
- bcrypt（或 argon2、scrypt）是**刻意設計得很慢**的雜湊，而且自動加入隨機 **salt**：同一個密碼每次雜湊結果都不同，攻擊者無法用預先算好的表（彩虹表）反查。

---

## 11.2 ★ `RandomToken() (string, error)` — 安全亂數

**規格**：產生 32 bytes 的亂數，以 `base64.RawURLEncoding` 編碼（43 個字元、可以放在網址中）。

**提示**：`crypto/rand` 的 `rand.Read(b)`。

**觀念**：`math/rand` 產生的是**可預測**的偽亂數，只要知道種子就能算出所有結果。用來產生 session ID、密碼重設連結、API 金鑰時，**一定要用 `crypto/rand`**。

---

## 11.3 ★★★ `Sign` / `Verify` — 簽章 token

這是 JWT（JSON Web Token）的簡化版，原理完全相同。

**格式**
```
base64url(JSON payload) + "." + base64url(HMAC-SHA256(payload 字串, secret))

例：eyJ1aWQiOjQyLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NjcyNzIwMDB9.kX9...
    └──────────── payload（任何人都能解碼看到內容）──────────┘ └ 簽章 ┘
```

**規格（`Verify`）**
- 格式錯誤、簽章不符、payload 無法解碼或解析 → `ErrInvalidToken`。
- `now` 已經到達或超過 `ExpiresAt` → `ErrExpired`。
- 測試會嘗試**竄改**：把 payload 改成 `role: admin`，沿用原本的簽章。

**提示**
- 編碼：`base64.RawURLEncoding`（URL 安全、無 `=` 補齊）。
- 簽章：
  ```go
  mac := hmac.New(sha256.New, secret)
  mac.Write([]byte(payload))
  sig := mac.Sum(nil)
  ```
- 比較簽章**一定要用 `hmac.Equal`**，不能用 `==` 或 `bytes.Equal`。
- 順序：**先驗證簽章，再解析 payload**。不要去處理未經驗證的資料。
- `strings.Cut(token, ".")` 拆成兩段；如果拆完 signature 中還有 `.`，簽章自然不會相符。

**觀念**
- **payload 沒有加密**，任何人都能用 base64 解碼看到內容。不要在 token 中放密碼、身分證字號等機密資料。
- 簽章的意義是**防竄改**：沒有 secret 就算不出正確簽章，所以伺服器可以信任 payload 的內容是自己簽發的。
- **時序攻擊**：`==` 比較字串時，遇到第一個不同的字元就停止。攻擊者可以量測回應時間，一個字元一個字元猜出正確的簽章。`hmac.Equal` 不論內容如何都花一樣的時間。
- 為什麼 `Verify` 要接收 `now` 參數而不是直接呼叫 `time.Now()`？這樣測試就能檢查「一小時後過期」，而不用真的等一小時——這叫**依賴注入**。

**延伸挑戰**
- 改用正式的 JWT 函式庫 [`github.com/golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)，比較它多處理了哪些事情（`alg` 標頭、`nbf`、`iat`……）。
- 實作 refresh token：access token 15 分鐘過期，refresh token 7 天、存在資料庫中，而且可以撤銷。

---

## 11.4 ★★ `RateLimiter` — 令牌桶限流

**規格**
- 每個 key（例如 IP）有自己的桶子，最多裝 `burst` 個令牌，一開始是滿的。
- 每秒補充 `rate` 個令牌（可以是小數，例如 0.5 = 每 2 秒 1 個），但不超過 `burst`。
- `Allow(key)`：桶子裡至少有 1 個令牌 → 拿走一個並回傳 `true`；否則 `false`。
- 必須可以安全地併發使用。

**提示**
- 不需要背景 goroutine 定時補充！每次 `Allow` 時再依照「距離上次經過多久」一次補足：
  ```go
  elapsed := now.Sub(b.last).Seconds()
  b.tokens = min(burst, b.tokens + elapsed*rate)
  b.last = now
  ```
- 令牌數用 `float64`，才能處理「0.5 秒補了半個」這種情況。
- 時鐘使用參數 `now func() time.Time`。測試會傳入一個假時鐘，直接「快轉」時間。

**觀念**：限流可以防止暴力破解密碼、API 被濫用、單一使用者把伺服器拖垮。正式專案可以用 `golang.org/x/time/rate`，演算法和你寫的相同。

**延伸挑戰**：map 會隨著不同 IP 無限增長（記憶體洩漏！）。加一個清理機制，移除超過 10 分鐘沒有使用的桶子。

---

## 11.5 ★★★ `Authenticate` / `ClaimsFrom` / `RequireRole` — 認證中介層

**規格**
- `Authenticate`：
  - 讀取 `Authorization: Bearer <token>` 標頭（注意 `Bearer` 後面有一個空白）。
  - 缺少標頭、格式不對、token 無效或過期 → **401**。
  - 驗證成功 → 把 `Claims` 放進 request 的 context，交給 `next`。
- `ClaimsFrom(ctx)`：取出 Claims；沒有時回傳 `false`。
- `RequireRole(role, next)`：沒有 Claims → **401**；角色不符 → **403**。

**提示**
- 放入與取出 context：
  ```go
  type ctxKey struct{} // 未匯出的型別：其他套件不可能用到同一個 key

  ctx := context.WithValue(r.Context(), ctxKey{}, claims)
  next.ServeHTTP(w, r.WithContext(ctx))

  c, ok := ctx.Value(ctxKey{}).(Claims)
  ```
- `strings.CutPrefix(header, "Bearer ")`。

**觀念**
- **401 Unauthorized**：「我不知道你是誰」（沒登入或 token 無效）。**403 Forbidden**：「我知道你是誰，但你不能做這件事」。
- 401 回應中**不要**透露詳細原因（「簽章錯誤」或「已過期」），以免提供攻擊者線索。詳細原因寫進伺服器日誌就好。

---

## 11.6 ★★ `RateLimit` — 限流中介層

**規格**：以客戶端 **IP** 為 key 呼叫 `l.Allow`，被拒絕時回應 **429 Too Many Requests**（建議加上 `Retry-After` 標頭）。同一個 IP 的不同 port 要算同一個人。

**提示**：`r.RemoteAddr` 的格式是 `"ip:port"`，用 `net.SplitHostPort` 拆開。

**思考題**：如果伺服器在 Nginx 或雲端負載平衡器後面，`RemoteAddr` 會變成負載平衡器的 IP，所有人共用同一個桶子。這時要改看 `X-Forwarded-For` 標頭——但這個標頭客戶端可以**任意偽造**。怎麼做才安全？

---

## 綜合延伸挑戰：完整的會員系統

結合練習 08、09、11，做一個有註冊與登入的 API：

| 端點 | 說明 |
| --- | --- |
| `POST /register` | `{email, password}` → 密碼雜湊後存入 users 表；email 重複回 409 |
| `POST /login` | 驗證密碼，回傳 `{token}`；套用嚴格的限流（每 IP 每分鐘 5 次） |
| `GET /me` | 需要登入，回傳自己的資料 |
| `GET /admin/users` | 需要 admin 角色 |

**安全檢查清單**
- [ ] 登入失敗時，「帳號不存在」和「密碼錯誤」回傳**完全相同**的訊息（避免被用來探測哪些 email 已註冊）
- [ ] secret 從環境變數讀取（`os.Getenv`），**不寫在程式碼裡、不 commit 進 git**
- [ ] 日誌中不記錄密碼與 token
- [ ] 正式環境只走 HTTPS
