## mutex

什麼是mutex?

講mutex之前我們需要先知道，儘管goroutine使用了channel來對資料進行同步
但是某些狀況下還是會需要避免競爭

例如兩個goroutin同時進行
```go
counter++
```
這並不是原子操作，他會被分解成(原子操作就是沒辦法分解成更小步驟的一個操作，天生帶了鎖的一個效果，其實就類似於硬體鎖)
```
load counter (100)
add 1 (101)
store counter (101)
```
如果兩個 core 同時做，會發生：

core1: load 100
core2: load 100
core1: store 101
core2: store 101

結果就少加了一次
這就是 race condition

那學java的時候基本上一定會碰到跟學到的就是CAS

```cpp
xchg 
compare-and-swap (CAS)
test-and-set
```
xchg 本身是原子操作，它象徵著把記憶體中的值，與某個暫存器做交換


compare-and-swap 


test-and-set 是 CPU 最原始的「互斥鎖基本指令」

他的邏輯是
```
old = memory
memory = 1
return old
```


cas的意思是
```java
if memory == expected {
    memory = newValue
    return true   // 說你成功
} else {
    return false  // 失敗，不修改
}
```
也就是說你會輸入一個參數，你認為當下這個數字是多少，接著就會知道成功與否，
所以儘管是無鎖操作但也常常會聽到CAS導致CPU耗能更高的問題

那接下來就可以進入Mutex的結構了
```go
type Mutex struct {
    state int32
    sema  uint32
}
```
go的mutex其實蠻簡單的就是有兩個欄位
state 用來表示鎖當前的狀態，用0跟1就可以達成
但如果只有0跟1 他就不需要int32這麼大了
他實際上還有一些參數

```
| starving | waiters | locked |
|   1bit   | 29bits  | 1bit   |
```
state 位元結構

```
bit 0     = locked
bit 1     = woken
bit 2–30  = waiter count
bit 31    = starving 
```
starving 用來表示 公平模式，也就是說讓等待最久的goroutine優先執行

那go整個運作流程就是
```
Lock()
CAS 嘗試搶鎖

成功 → 拿到鎖

失敗 → 進入自旋

自旋（spin）數次

適合在多核心環境

若有機會奪鎖 → 快速完成

仍無法搶到鎖 → 進入阻塞

waiter count++

goroutine 進入 sema 中 wait

Unlock()

清除 locked

若有等待 → 喚醒一個 goroutine

```

最後就來示範一些mutex的使用

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex
    counter := 0

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        mu.Lock()
        counter++
        mu.Unlock()
    }()

    go func() {
        defer wg.Done()
        mu.Lock()
        counter++
        mu.Unlock()
    }()

    wg.Wait()
    fmt.Println("counter =", counter)
}
```
上面可以看到一個很簡單粗暴的加鎖範例

Mutex 常用來保護“資料”，所以慣例是把它放在 struct 裡
```go
type SafeCounter struct {
    mu sync.Mutex
    v  int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    c.v++
    c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.v
}

func main() {
    c := &SafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Inc()
        }()
    }

    wg.Wait()
    fmt.Println(c.Value())
}

```

