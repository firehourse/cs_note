## receiver


為了瞭解go 的整個error 運作，我們會從receiver開始了解，

其實receiver不是太難的東西，我一開始問AI問了很久走了很多彎路
其實用C\C++或者Asemble的角度來看就非常的簡單了

首先我們宣告一個簡單的dog結構體，不需要塞任何東西給他
```go
type dog struct {}
```
接著我們宣告一個方法

```go
func (dog) bark() {
    fmt.Println("bark")
}
```
這時候應該會有點疑惑這個dog 跟 bark 有什麼關係，他們怎麼被綁訂了

其實就很簡單，不管是哪個語言，如果你看過他們的Asemble code就會知道
這個dog 實際上就類似於其他語言的this或者self

但go 的receiver 運作的方式其實跟cpp完全不同

但是在cpp裡面，this永遠是一個指標指向了當前物件
go這邊比較tricky，雖然妳看起來是 a.add()這種調用
但其實go會變成add(a)這樣的方式，那既然是這樣的調用他當然就可以選擇要傳指標還是value了
recevier只是做了一個型別檢查的功能，來達到讓方法聚合的功能

```go
func (d dog) bark() {
    fmt.Println("bark")
}
```
這時候的d 就是傳值

```go
func (d *dog) bark() {
    fmt.Println("bark")
}
```
這時候的d 就是傳pointer

那這有什麼好處呢？

舉一個最近的例子
我有一個DTO他要進行傳輸，傳輸的途中我需要對這個DTO的值進行填入其他的DTO給不同的service使用
但當我做某些操作的時候一不小心就改變了他的值
這會導致在後續的流程中出問題，在OOP的世界中
我就需要非常麻煩的把它寫成getter跟setter來保證我取值出來進行操作的時候是clone的

在go的世界中
我只需要
```go
type dto struct {
    uint64 x
    uint64 y
}

func (d dto) add(x uint64, y uint64) uint64 {
    return d.x + d.y
}

d := dto{1,2}
// 這時候dto 的值是 1 跟 2 
p := d.add(10,20)
// 接著我試著打印他們

fmt.Println(d)// 1 2
fmt.Println(p) // 3
```

然後接著是一個非immutable的例子
```go
func (d *dto) add(x uint64, y uint64) uint64 {
    d.x = x
    d.y = y
    return d.x + d.y
}
d := dto{1,2}
// 這時候dto 的值是 1 跟 2 
p := d.add(10,20)
// 接著我試著打印他們
fmt.Println(d)// 10 20
fmt.Println(p) // 30
```

這樣子我們就完成了receiver的介紹
