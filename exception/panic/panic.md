## panic

什麼是panic?

panic 就是程式運行遇到了無法處理的狀況，那當突發狀況發生的時候，代碼自然無法運行下去，可以說panic也是一種針對異常的處理方式
只是代價比較昂貴而已，但他避免了更多的意外狀況

Go 的思想 他並不把panic當作是 錯誤 而當作是
一種不正常的工作流
其實就像這次cloudflare的事件，是由unwrap引發的，而這也是panic 的處理方式
那他的處理方式是非常重要的，也就是讓runtime執行stack unwind
為什麼說這非常重要呢?
要是你不進行stack unwind，你可能就有一塊記憶體空間被永遠的佔用除非你斷電重開不然就不見了

當程式panic 他會做什麼?

1. 建立 panic value（可以是任何型別）
2. 開始 stack unwind（從最深一層往外爆）
3. 每到一層，就執行該層的 defer
4. 若遇到 recover → 中止 panic，回到正常路線
5. 若沒有 recover → 程式 crash（印 stack trace）

那什麼是panic value?

panic value 可以是任何型別，但通常我們會用error interface來包裝
```go
func panic(v interface{}) { ... }
```
所以 panic value = interface{}
→ 存在 g.panicVal 裡
→ recover() 也會拿到它

看runtime code你就會看到
```go
panic(gp, val)
```
那什麼時候可以用panic?

例如我們在 php 的時候常常會寫switch case 然後在某些特殊情況我們也許會返回null object或者給個空陣列之類的
這當然是一個處理方法，但他可能會導致我們的數據非常的髒，但同時也讓程式崩潰這件事情被隱藏了起來
但這時候我們可以用一個panic來處理，這樣至少在測試的時候有機會知道這會發生錯誤而在上線的時候可以配合zabbix或者error log 來進行監控

```go
switch v.(type) {
case int:
case string:
default:
    panic("unexpected type")
}
```
那最後就是 panic 這東西其實就近似於其他語言的exception ，他主要還是由 go runtime 去建造出來的一個panic chain 
那我們在宣告panic value的時候其實就相當於向裡面某一層添加一個arg而已
