## defer

defer 是go 特有的一個語法，他類似於cpp的 destructor
defer 是compiler + runtime 合作生成的一段 延遲執行的函式

他會在function return 前執行
在 panic stack unwind 的時候也會執行

### defer 是怎麼被compiler生成的

先看一段簡單的code
```go
func f(){
    defer fmt.Println("A")
    defer fmt.Println("B")
}
```

Go compiler 實際上做的是：
把每個defer 轉成一個runtime呼叫
同時建立一個 __defer
把defer 用link list 的方式串起來

runtime的 __defer 結構可以簡化為

```go
type __defer struct {
    next    *_defer      // 指向上一個 defer（LIFO）
    fn      *funcval     // 要執行的函式
    pc      uintptr      // 呼叫點位址
    sp      uintptr      // stack pointer
}
```
那我們前面呼叫了 A 跟 B 他就會變成
defer1 := &_defer{ fn: A }
defer2 := &_defer{ fn: B, next: defer1 }

那這其實就跟一般 stack 的 LIFO 一模一樣

這讓defer的執行順序可以跟一般的function 呼叫一致，也就是當你呼叫鏈到最後的時候接著就會一個個pop出來

那程式在return 的時候就會開始釋放defer 有沒有覺得跟javascript 的 async/await 有點類似，不過async/await是 queue 的方式

而當程式發生panic 的時候，go的runtime會開始釋放defer，所以就可以利用這個特行來攔截或對panic進行處理

那這邊我們就需要說到一個很重要的概念，程式的崩潰並不意味著那個進程或者go的runtime更或者電腦崩潰了
實際上程式的崩潰是觸發了意想不到的例外狀況並且沒有被妥善處理
所以panic 並不是一個錯誤，而是一個特殊的狀況已經被知悉，但返回了這個問題你沒有去處理
真正會讓系統掛掉的反而是memory leak、kernel panic 這類OS層級的錯誤問題

簡單的觸發一個panuc就可以驗證我所說的
```go
func f() {
    defer fmt.Println("A")
    panic("boom")
    fmt.Println("never run")
}
```
執行結果會是
```
A
panic: boom
<stack trace>
```
雖然panic 會讓程式掛掉，但defer 仍然會被執行

defer 是唯一能攔住 panic 的機制，recover 一定要在 defer 裡