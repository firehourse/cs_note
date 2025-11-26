## error chain
這邊準備要來談談error chain

那因為php 算是一個簡單易讀的語言，所以我會先用php 來手刻一下

雖然php的interface跟 go的差異很大，但概念上應該可以相似畢竟現在也都使用IoC的概念居多
```php
interface ErrorInterface{
    public function message(): string;
}

# 接著我要實作他

class SimpleError implements ErrorInterface{
    private $message;
    public function __construct($message){
        $this->message = $message;
    }
    public function message(): string{
        return $this->message;
    }
}

# 接著演示一下怎麼使用他

// 這邊我宣告了一個foo ，他的回傳必須要是一個ErrorInterface或者null
function foo() : ?ErrorInterface{
    return new SimpleError("hi");
}

err = foo();
if ($err !== null){
    echo $err->message();
}
```
從上面可以看到，他的概念其實很簡單就是當你在返回的時候啊，無論如何都可以用一個errorInterface來包裝，但go 比較特別的點是他可以接受多值返回
所以通常我們會讓返回結果一個是原本要return的一個則是error

那當然Rust 的 Result 也是乘載了這種精神，但他把成功與失敗都包在了一個Result struct 裡面，這件事情其實我們也可以手刻一個出來，但我想這就留給一個獨自的篇章

但這邊我會先試著用Go 來寫一個error -> error chain 這樣的形式

那複習一下我們提過的Error interface
```go
type error interface {
	Error() string
}
```
這時候做一個SimpleError

```go
type SimpleError string
// 這邊是實現 error interface
func (e SimpleError) Error() string {
    return string(e)
}
```

那這簡單的error 就完成了，最後是也很簡單的eroor chain
那名字是chain所以很快的就可以聯想到link list 結構
我們會把錯誤一層一層的用一個Wrapper包起來

```go
// 先創建一個包裝器資料結構
type ErrorWrapper struct {
    msg string
    cause error
}
// 我們要實踐他的方法，其實非常類似責任鏈的set handler

func (e *ErrorWrapper) Error() string {
    // 避免空指針
    if e.cause == nil {
        return e.msg
    }
    return e.msg + ": " + e.cause.Error()
}
// 這時候已經完成一個error chain的結構了，那如同鏈表的讀取，他需要一層層被遍歷拆開，這樣我們就可以追蹤底層的錯誤

func (e *ErrorWrapper) Unwrap() error {
    return e.cause
}
```
這基本上就是 result pattern 的一個實現，每個語言其實都可以實作出來，只是go一開始就在語法層級兼融了這樣的模式