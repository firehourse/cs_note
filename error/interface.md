## interface

想要搞懂 go 的 error chain 實作繞不開的就是 go 的 interface

跟java、php 這類的interface用法不同，以目前的大宗用法，interface更常被IoC解耦使用
那就算是原生的interface，也跟go的interface使用上有差異

在php與java中，interface會顯示的宣告他的function，然後你的類別需要去implements 宣告實作這個interface來達成綁定

在go中，interface是隱含的，也就是go的compiler會自動幫你包裝
只要你符合interface的定義，你就是這個interface的實作者，interface只是一個寬鬆的對於某些function的定義

也就是說當你這個function實作了 interface裡面的function，那他就被視為符合interface的定義

那說回到error interface這件事情

```go
type error interface {
    Error() string
}
```
看到上面的代碼，我們宣告了一個error interface而它裡面只有一個Error function
這代表當任何時候我們只要宣告某個struct符合error interface的定義，我們就可以說他是一個error implementation

當我們宣告
```go
type MyErr struct{}
func (MyErr) Error() string { return "hi" }

var e error = MyErr{}
```
這邊在compile 階段，compiler會判定MyErr是否符合error interface的定義
如果符合，compiler就會幫我們包裝，
但建立 interface value（裝進 e） 是 runtime 做的

Go runtime 會建立interface value 結構
```go
interface value
    itab pointer -> 指向了 interface+具體類型的方法表
    data pointer -> 指向了具體類型的值
```
簡化成
```go
e = {
    tab: &itab{ type: MyErr, interface: error, methods: ... }
    data: pointer to MyErr{}
}
```
其實就可以看做 c的vtable java的 class meta data 或者 php的zend_class_entry 還有 python 的 __dict__
