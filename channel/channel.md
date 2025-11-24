# channel

channel 是 go 特有的一個概念，主要是在於多線程或者協程之間，你還是必須要用一個鎖來保護資料的正確性

我想很多人應該都有這個概念就是，在現在多執行緒的時代，你沒辦法保證一筆資料他被多個執行緒同時存取的順序

用一個經典的python 例子

```python
x = 0

def add():
    global x
    for _ in range(500):
        x += 1  
t1 = Thread(target=add)
t2 = Thread(target=add)

t1.start()
t2.start()

t1.join()
t2.join()

print(x)
```
這裡你會預設操作出來的結果應該會是1000，但實際上你會得到一個不確定的結果

那 go 的思想是，他不在資料上面加上一個lock 來保證數據的正確性
而是用一個通信的queue 來保證資料的正確性，那lock的操作就會體現在這個queue上面
那麼channel 本身只是一個queue，但他需要滿足一些條件才能被觸發
當 goroutine 對 channel 執行 send/recv 時，如果條件不滿足，
goroutine 會被 runtime 調度器 park（阻塞）
直到另外一個 goroutine 配對（send 或 recv），才會 unpark（喚醒）。
這種阻塞/喚醒行為就是造成 channel 行為的關鍵，也會帶來一些設計上的問題。

那因為我比較不擅長抽象的理解，只能舉出一個簡化的channel

```go
type hchan struct {
    qcount   uint           // buffer 中目前的元素數量
    dataqsiz uint           // buffer 的大小（0 = 無 buffer channel）
    buf      unsafe.Pointer // 指向 buffer（環狀 queue）

    elemsize uint16         // 每個元素的大小
    closed   uint32         // channel 是否被 close
    
    sendx    uint           // 下次送值的 buffer index
    recvx    uint           // 下次收值的 buffer index

    recvq    waitq          // 等待接收的 goroutine queue (G 列表)
    sendq    waitq          // 等待發送的 goroutine queue (G 列表)

    lock     mutex          // 保護整個 hchan
}

```
這邊請gpt畫了一個簡單的圖

```
      +----------------------+
      |       hchan          |
      +----------------------+
      | qcount     = 2       |
      | dataqsiz   = 5       |
      | buf  ---> [A][B][ ][ ][ ]  (環狀 queue) 
      | sendx      = 2       |
      | recvx      = 0       |
      | closed     = 0       |
      |
      | recvq --->  G3 -> G7 -> G9    (阻塞等待 recv)
      | sendq --->  G2 -> G4          (阻塞等待 send)
      +----------------------+
```

那這上面其實還有點不清楚，首先我們要理解的是buffer

```go
ch := make(chan int)
```
這個channel 緩衝大小為 0 
也就是說送資料時必須要馬上送，收資料時也必須要馬上收
否則就會阻塞

```go
ch <- 1 // 送資料
```

```go
x := <- ch // 收資料
```
這種channel 叫做 rendezvous channel
他不能算是queuue而是handshake channel

而有buffer的情況下
```go
ch := make(chan int, 3)
```
這個channel 他就會產生
[] [] []三個這樣的空間
你就可以連續的送入
```go
ch <- 1
ch <- 2
ch <- 3
```
buffer 變成
```go
[1][2][3]
```
而這種狀況的channel 只有在滿載與全空的情況下會阻塞

那我覺得要更實際的講整個channel 的運作就必須要談到GMP模型的調度了

那如果不熟悉thread與process這邊簡略介紹一下
process就是一塊記憶體的劃分 底下最少會有一個os thread工作，而一個thread則是執行的最小單位
但是協程可以視為一個更迷你的單位，他實際上並不為os所知道，os只知道他調度了一個thread
而我們則是用資料結構做出更迷你的一個儲存單位讓thread可以切換，也就是說goroutine或者說協程
可以視為超迷你的process，但他可以在一個公共的process下運作

假設現在有
1個M
1個P
多個G

只有一個M1的情況下，有兩個G使用了channel進行溝通
G1: 發送 ch <- 1

G2: 接收 x := <-ch

G1會將資料送入channel
G2會等待channel有資料

當channel有資料時
G2會將資料取出
G1會繼續

runtime 對 G1 做：

把 G1 放到 hchan.sendq

G1 狀態變成 waiting

暫停（park）G1

scheduler 接管 M1，換別的 G 執行

那能夠這樣有一個很重要的原因是 GO 強制了 Goroutine的 上下文切換時間

以前（Go 1.1～1.13）主要是 cooperative（依靠函式呼叫 safepoint 才能切換 G）。

但後來（Go 1.14 之後）：

Go runtime 會強制在 goroutine 執行太久時，把它搶下來切換到別的 goroutine。

也就是說：

✔ 即使只有「1 個 M」，scheduler 也能隨時搶佔正在執行的 G
✔ 即使某個 G「沒有主動讓出 CPU」，runtime 也會定期中斷它
✔ 因此 G1/G2 都有機會被調度，channel 才能配對 send/recv


那也就是說，就如同kafka的高吞吐設計，當你給channel設置了緩衝區，Thread會不停的切換goroutine來消化你這些資料
這樣想應該就清楚了很多
那我們來看一個例子


```go
    func main(){
        c1 := make(chan string)
        c2 := make(chan string)
        // 上面宣告了兩個channel
        go func(){
            for {
                c1 <- "from 1"
                time.Sleep(time.Second * 1)
            }
            }()
        go func(){
            for {
                c2 <- "from 2"
                time.Sleep(time.Second * 2)
            }
        }()
        for {
            fmt.Println(<-c1)
            fmt.Println(<-c2)
        }
    }
```
那這邊的運行就是說，當我們宣告了goroutine 他就開始運行了，那前面提到沒有buffer的情況下channel send 就要等待 recv
所以這邊 go 1送出來 然後等待1秒接著我們會收到c1的資料 然後收到c2的資料
但你會發現一個狀況，為什麼他們都是一起打印出來的呢？實際上是因為這邊雖然channel 的 send我們用兩個goroutine處理了
但是我們的recv卻是用一個goroutine來處理，也就是main

所以這時候就輪到select語句出場了
```go
for {
    select {
    case msg := <-c1:
        fmt.Println(msg)
    case msg := <-c2:
        fmt.Println(msg)
    }
}
```
那select 語句具體做了什麼呢?
他其實就是，當誰先送到了那誰就執行，注意，他不是用兩個goroutine來處理，
而是像我們一般的switch case一樣 誰符合條件誰就觸發，而他的唯一條件就是，這個case不阻塞他就會觸發