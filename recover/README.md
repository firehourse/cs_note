# recover

recover 的本質實際上不是拿來捕捉 panic 的
他實際上是終止 stack unwind 的一個機制

他必須要被放在 defer 中，因為當panic發生的時候
只有defer 會在 stack unwind 的過程中被執行

那對於recover我覺得比較需要提到的是
他終止了stack unwind那他是怎麼恢復程序的執行的?

前面有提到panic 也是藉由一個chain開始進行返回，

recover 做的事情其實就三步：

1. 取 panic value（p.arg）
2. 標記 p.recovered = true
3. 從 panic chain 中移除這個 panic（pop）

也就是

panic2, panic1 → 被恢復
panic3 → 被 recover →結束 unwind

來看看例子
```go
func demo() {
    defer func() {
        if r := recover(); r != nil {
        // 當panuc 發生的時候recover就不會是nil的，就跟error interface一樣
            fmt.Println("recovered:", r)
        }
    }()
    panic("boom")
    fmt.Println("unreachable")
}
```
整個流程就會是

panic("boom")
 ↓
runtime.gopanic 建立 panic frame
 ↓
unwind 1 層 → 遇到 defer
 ↓
執行 defer → recover() → 截斷 unwind
 ↓
回到 demo() → demo() 正常 return


那何謂正常return呢? 其實就如同C語言 return 0，如果你沒有在recover裡面特別處理它就會當作沒事直接的return然後繼續執行該defer接下來的代碼

看完了流程整個recover唯一要注意的應該就是要放在stack unwind的流程裡面，也就是說panic發生開始釋放的時候要能遇到這個Recover不然形同虛設


go 社區建議將recover 用在以下兩種場景，我覺得這也是其他語言可以借鑑的

在最頂層的時候進行保護，並且將panic轉換為error並記錄log
避免整個服務崩潰

當寫某些內部的邏輯的時候，有可能有一些預期外的輸入，這時候就可以捕獲異常來返回，其實就近似於php常用的try catch