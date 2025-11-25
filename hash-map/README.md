接下來要講的其實是go的hashmap，那我其實沒手寫過hashmap，不過有幸看過java 對 string 做 hash 的 source code，所以這邊嘗試用go簡單的刻一下hash的部分
```go
func stringHash(s string) int {
    h := 0
    for i := 0; i < len(s); i++ {
        h = 31 * h + int(s[i])
    }
    return h
}
```
實際上這個hash 只是說會算出一串數字，然後我們會取他的餘數
```go
func bucketIndex(h int){
    return h % 16
}
```
這時候真正的得到了一個位置，會將要放入的string 丟到一個由go runtime 管理的 array 來放
```go
func put(key string, value string){
    h := stringHash(key)
    index := bucketIndex(h)
    bucket[index] = append(bucket[index], key, value)
}
```
那下次要取，我們只要經過一樣的計算，就可以到同樣的bucket去取，藉由這樣的方式來達到直接訪問的效果


# 常見的幾個 go 宣告 map 的方法

最基本的，用make來創建一個map
```go
var m map[string]int
```
map 的零值是 nil，不能直接寫入，會 panic，所以正常都要用以下方式創建

```go
m := make(map[string]int)
```
基本上就是 map[<key 的型別>]<value 的型別> 這樣的宣告方式

也可以直接給一個初始值
```go
m := map[string]int{"foo": 1, "bar": 2}
```

這樣的方式，就是直接取出map中的值
```go
val := m["foo"]
```

這樣的方式，就是直接將map中的值改掉
```go
m["foo"] = 3
```



那go 有個地方跟 cpp 蠻像的，就是判斷一個key是否存在
因為go有一個特性是如果key不存在他會返回0而不會像php undefind index
所以一定要注意要判斷key是否存在

```go
value, ok := m["apple"]
if ok {
    fmt.Println("exist:", value)
} else {
    fmt.Println("not exist")
}
```

Go 可以這樣直接刪除map中的值，即使key不存在也不會報panic
```go
delete(m, "foo")
```


在php 中最常見的 foreach 陣列，這邊有一點我不是很確定的是這是不是也是靠iterator的形式來達成的，畢竟他也是一個make出來的struct
```go
for k, v := range m {
    fmt.Println(k, v)
}
```

go 可以通過下面的方式輕易地取得 map 的長度
```go
n := len(m)

```

map的value可以是任何東西，類似於cpp的vector，非常的自由
```go
nested := make(map[string]map[string]int)

nested["outer"] = map[string]int{
    "inner": 1,
}
```

那比較特別的是 map 像php的物件，他是reference type的 所以使用上要注意
```go 
m1 := map[string]int{"a": 1}
m2 := m1
m2["a"] = 100
fmt.Println(m1["a"]) // 100 !!!
```
這樣的話，他會直接改掉m1的值，因為他是reference type的

那他也不像slice有提供copy或者其他語言的clone語法，需要自己寫一個func 來做複製
```go
func cloneMap[K comparable, V any](src map[K]V) map[K]V {
    dst := make(map[K]V, len(src))
    for k, v := range src {
        dst[k] = v
    }
    return dst
}
```

go的 map key是無序的，有點類似python 的 set ，所以每次印出可能順序都會不同
另外就是跟 java 的 hashmap一樣的老問題，本身是不支援concurency的 java 採用了ConcurrentHashMap 那go 是使用sync.map不過原理我還不是很確定，有空再來詳細看一下