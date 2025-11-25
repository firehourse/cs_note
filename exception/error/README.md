這邊來聊聊 Go 的 error。
其實我過去也沒有很深入研究過 Go 的 error system，只知道有些人常拿 Rust 的 Result 來跟 Go 比。有個 up 主還說 Rust 的 Result 是借鑒了 Go 的 error（雖然最近 Cloudflare 因為 unwrap() 出事確實有點好笑）。

首先 go 的 error 是需要你顯式去掉用的，也就是說在function 宣告的時候就決定了return會塞一個 error
```go
func DoSomething() (Result, error)
```
你必須自己寫：
```go
v, err := DoSomething()
if err != nil {
    return nil, err
}
```

那這個error的結構非常的簡單

```go
type error interface {
	Error() string
}
```
但這底下蘊含的東西可就多了，首先他是個interface，但 interface 本身不是這個 struct，而是「包裝」

每個 interface 在 runtime 都會被包裝成以下結構：
```go
type iface struct {
    tab  *itab          // 指向型別資訊（包含方法表）
    data unsafe.Pointer // 指向實際資料（例如 errorString）
}
```

那當你寫下
```go
var err error = errors.New("hi")
```
這時候發生的事情是
先創建 struct：
```go
errorString{s: "hi"}
```
接著這個 struct 被「轉成 interface」
```go
err = iface{
    tab:  →  errorString 的 type info（itab）
    data: →  指向 errorString{s: "hi"} 的記憶體位置
}
```
所以平常 cdp, err := DoSomething() 的 err 就需要判斷他有沒有接到東西，如果有就是有接到錯誤那你就需要去排除

簡單的error處理大概就是這樣，接下來我會在下一篇講到error真正的鏈式結構，這才是他被說跟result像的地方（其實是我自己還沒吃透）

如果想要深入的了解go 的error 我會建議從我的receiver篇開始 receiver -> interface -> error_chain 這樣的順序來理解