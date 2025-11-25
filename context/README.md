# context

什麼是context? 

context.Context 是 Go 語言內建的一個 控制請求生命週期的機制，用來在 goroutine 之間傳遞：

取消信號（cancel）

timeout / deadline

每個請求的 metadata（key-value）

它本質是一組 interface：
```go

type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}

```
那既然有interface他就會有實作

go 語言不希望每個人都實作自己的context，所以他提供了預設的context
概念上偏向單例模式
你會去import context的 package

而裡面有4個核心的struct
emptyCtx        Background(), TODO()
cancelCtx       WithCancel()
timerCtx        WithDeadline(), WithTimeout()
valueCtx        WithValue()

官方禁止你 new，就是因為：

若每個人都亂 new，就破壞了 context 的「沿著請求生命週期、可取消、可溯源、樹狀傳播」的設計。

## cancelCtx
```go
type cancelCtx struct {
    Context               // parent
    mu       sync.Mutex
    done     chan struct{}
    children map[canceler]struct{}
    err      error
}
```
其中done 就是一個channel，他用來廣播cancel的訊號，還記得前面提過go 使用channel來溝通嗎?
這裡使用上其實就是我們會把這個信息廣播出去給所有的goroutine去接收，當他得知這個訊號時就會停止

children則是一個被維護的陣列，他記錄了我們針對cancelCtx所產生的子context
mu 則是Mutex，可以想像成lock，用來保護children的維護
err 就是之前的error interface 沒什麼特別的

到這邊可以發現其實我們都是用基本的資料結構工具把這些struct組裝起來

## timerCtx 
timerCtx實際上就是cancelCtx的上層封裝，他會多一個timer
在時間到的時候呼叫cancel，就這麼樸實無華
```go
type timerCtx struct {
    cancelCtx
    timer *time.Timer
}
```

## valueCtx
valueCtx也很簡單，他其實也是必須要在裡面塞一個context，然後塞入key value 
```go
type valueCtx struct {
    Context
    key, val any
}
```
實際使用上大概像這樣
```go
ctx = context.WithValue(ctx, "uid", 123)
```
那其實就像是繼承關係，但實際是一個link list 也就是說 我呼叫這個ctx的時候，它裡面包了一個context可以指到更上層的context
當你要找key的時候找到就回傳，爬到最後就結束，就是這麼的簡單

## emptyCtx
這應該是最單純的部分了，他其實就是作為一個根結點
所以結構就只有
```go
type emptyCtx int
```
最後就是提一下，當你呼叫了context的 cancel的時候

由於當下只會釋放當前ctx的資源，但是子ctx之間還是會層層的被引用導致gc沒有被回收
所以必須要在defer中再額外呼叫cancel

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 使用 ctx
}
```
cancel 詳細會做的事很多，之後也許可以另開一篇來深度剖析